package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	KafkaServer = "localhost:9092,localhost:9093,localhost:9094"
	KafkaTopic  = "inventory-ticks"
)

var producer *kafka.Producer
var err error

func InitKafkaProducer() {
	producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
		"client.id":         "myProducer",
		"acks":              "all",
	})
	if err != nil {
		panic(err)
	}

	// configure the producer to produce data asynchronously
	producer.ProduceChannel()

	go func() {
		for event := range producer.Events() {
			switch message := event.(type) {
			case *kafka.Message:
				if message.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", message.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset)
				}
			}
		}
	}()

	// Handle signals to gracefully shut down the producer
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle signals
	go func() {
		<-signals
		fmt.Println("Received interrupt signal. Closing producer...")
		// Flush the producer before closing
		// waiting 15 seconds
		producer.Flush(15 * 1000)
		producer.Close()
	}()
}

func SendMessageToKafka(message []byte, topicName string) {
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)

	if err != nil {
		panic(err)
	}
}
