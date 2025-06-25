package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	KafkaServer      = "localhost:9092,localhost:9093,localhost:9094"
	KafkaTopic       = "inventory-ticks"
	MIN_COMMIT_COUNT = 1
)

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
		"group.id":          "group1",
		"auto.offset.reset": "earliest"})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{KafkaTopic}, nil)

	msg_count := 0
	run := true
	for run == true {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			msg_count += 1
			if msg_count%MIN_COMMIT_COUNT == 0 {
				consumer.Commit()
			}
			fmt.Printf("%% Message on %s:\n%s\n",
				e.TopicPartition, string(e.Value))

		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			fmt.Printf("Ignored %v\n", e) // application-specific processing
		}
	}

}
