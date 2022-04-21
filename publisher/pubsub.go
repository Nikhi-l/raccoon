package publisher

import (
	"context"
	"fmt"
	"log"
	"cloud.google.com/go/pubsub"
	"sync"
	"sync/atomic"
)

func NewPubSub(ctx context.Context)(*PubSub,error) {
	projectID:="test"
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
		return nil, err
	}
	return &PubSub{
		cli: pubsub.Client(*client),
	},nil
}

type PubSubPublisher interface {
	// ProduceBulk message to kafka. Block until all messages are sent. Return array of error. Order is not guaranteed.
	ProduceBulk(request collection.CollectRequest) error
}
type PubSub struct {
	cli           pubsub.Client
}
func (pr *PubSub) ProduceBulk(request collection.CollectRequest, ctx context.Context) error {
	events := request.GetEvents()
	errors := make([]error, len(events))
	totalProcessed := 0
	for order, event := range events {
		topic :=  event.Type
		var wg sync.WaitGroup
		var totalErrors uint64
		t := pr.cli.Topic(topic)
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("Value" + event.EventBytes),
		})
		wg.Add(1)
		go func( res *pubsub.PublishResult) {
			defer wg.Done()
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				fmt.Printf( "Failed to publish: %v", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Printf( "Published message; msg ID: %v\n", id)
		}(result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d messages did not publish successfully", totalErrors)
	}
	return nil

	}



