package app

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"raccoon/collection"
	"raccoon/config"
	"raccoon/services"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/publisher"
	"raccoon/worker"

	"runtime"
	"syscall"
	"time"
)

// StartServer starts the server
func StartServer(ctx context.Context, cancel context.CancelFunc) {
	bufferChannel := make(chan collection.CollectRequest, config.Worker.ChannelSize)
	httpServices := services.Create(bufferChannel)
	logger.Info("Start Server -->")
	httpServices.Start(ctx, cancel)
	logger.Info("Start pubsub publisher -->")
	PubSubPublisher, err := publisher.NewPubSub(ctx)
	if err != nil {
		logger.Error("Error connecting to Pubsub", err)
		logger.Info("Exiting server")
		os.Exit(0)
	}

	logger.Info("Start worker -->")
	workerPool := worker.CreateWorkerPool(config.Worker.WorkersPoolSize, bufferChannel, config.Worker.DeliveryChannelSize, PubSubPublisher)
	workerPool.StartWorkers()
	//go kPublisher.ReportStats()
	go reportProcMetrics()
	go shutDownServer(ctx, cancel, httpServices, bufferChannel, workerPool, PubSubPublisher)
}

func shutDownServer(ctx context.Context, cancel context.CancelFunc, httpServices services.Services, bufferChannel chan collection.CollectRequest, workerPool *worker.Pool, pubsub *pubsub.Client) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[App.Server] Received a signal %s", sig))
			httpServices.Shutdown(ctx)
			logger.Info("Server shutdown all the listeners")
			timedOut := workerPool.FlushWithTimeOut(config.Worker.WorkerFlushTimeout)
			if timedOut {
				logger.Info(fmt.Sprintf("WorkerPool flush timedout %t", timedOut))
			}
			//flushInterval := config.PublisherKafka.FlushInterval
			logger.Info("Closing PubSub client")
			//logger.Info(fmt.Sprintf("Wait %d ms for all messages to be delivered", flushInterval))
			err := pubsub.Close()
			if err!=nil{
				log.Fatalf("Could not close pubsub Client: %v", err)
			}
			eventsInProducer:=0
			/**
			@TODO - should compute the actual no., of events per batch and therefore the total. We can do this only when we close all the active connections
			Until then we fall back to approximation */
			eventsInChannel := len(bufferChannel) * 7
			logger.Info(fmt.Sprintf("Outstanding unprocessed events in the channel, data lost ~ (No batches %d * 5 events) = ~%d", len(bufferChannel), eventsInChannel))
			metrics.Count("kafka_messages_delivered_total", eventsInChannel+eventsInProducer, "success=false")
			logger.Info("Exiting server")
			os.Exit(0)
		default:
			logger.Info(fmt.Sprintf("[App.Server] Received a unexpected signal %s", sig))
		}
	}
}

func reportProcMetrics() {
	t := time.Tick(config.MetricStatsd.FlushPeriodMs)
	m := &runtime.MemStats{}
	for {
		<-t
		metrics.Gauge("server_go_routines_count_current", runtime.NumGoroutine(), "")

		runtime.ReadMemStats(m)
		metrics.Gauge("server_mem_heap_alloc_bytes_current", m.HeapAlloc, "")
		metrics.Gauge("server_mem_heap_inuse_bytes_current", m.HeapInuse, "")
		metrics.Gauge("server_mem_heap_objects_total_current", m.HeapObjects, "")
		metrics.Gauge("server_mem_stack_inuse_bytes_current", m.StackInuse, "")
		metrics.Gauge("server_mem_gc_triggered_current", m.LastGC/1000, "")
		metrics.Gauge("server_mem_gc_pauseNs_current", m.PauseNs[(m.NumGC+255)%256]/1000, "")
		metrics.Gauge("server_mem_gc_count_current", m.NumGC, "")
		metrics.Gauge("server_mem_gc_pauseTotalNs_current", m.PauseTotalNs, "")
	}
}
