package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	viper.AutomaticEnv()
	assert.Equal(t, "debug", LogLevel())
}

func TestAppPort(t *testing.T) {
	os.Setenv("APP_PORT", "8080")
	viper.AutomaticEnv()
	assert.Equal(t, "8080", AppPort())
}

func TestNewKafkaConfig(t *testing.T) {
	os.Setenv("KAFKA_BROKER_LIST", "kafka:6668")
	os.Setenv("KAFKA_TOPIC", "test1")
	os.Setenv("KAFKA_ACKS", "1")
	os.Setenv("KAFKA_QUEUE_SIZE", "10000")
	os.Setenv("KAFKA_FLUSH_INTERVAL", "1000")
	os.Setenv("KAFKA_RETRIES", "2")
	os.Setenv("KAFKA_RETRY_BACKOFF_MS", "100")

	expectedKafkaConfig := KafkaConfig{
		brokerList:     "kafka:6668",
		topic:          "test1",
		acks:           1,
		maxQueueSize:   10000,
		flushInterval:  1000,
		retries:        2,
		retryBackoffMs: 100,
	}

	kafkaConfig := NewKafkaConfig()
	viper.AutomaticEnv()
	assert.Equal(t, expectedKafkaConfig, kafkaConfig)
}
