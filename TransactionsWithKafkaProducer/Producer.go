package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	KafkaServer     = "localhost:9092,localhost:9093,localhost:9094"
	KafkaTopic1     = "inventory-ticks"
	KafkaTopic2     = "order-ticks"
	TransactionalId = "my-transactional-producer"
)

type ProductInventory struct {
	Id       int
	Name     string
	Quantity int
}

type ProductOrder struct {
	Id       int
	Name     string
	Quantity int
}

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
		"client.id":         "myProducer",
		"acks":              "all",
		"transactional.id":  TransactionalId,
	})
	if err != nil {
		panic(err)
	}

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

	topic1 := KafkaTopic1
	topic2 := KafkaTopic2
	productName := "productName_1"
	productInv := ProductInventory{
		Id:       1,
		Name:     productName,
		Quantity: 5,
	}
	value1, err := json.Marshal(productInv)
	if err != nil {
		panic(err)
	}

	productOrder := ProductOrder{
		Id:       1,
		Name:     productName,
		Quantity: 3,
	}
	value2, err := json.Marshal(productOrder)
	if err != nil {
		panic(err)
	}

	maxDuration, err := time.ParseDuration("10s")
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), maxDuration)
	defer cancel()

	err = producer.InitTransactions(ctx)
	if err != nil {
		panic(err)
	}

	err = producer.BeginTransaction()
	if err != nil {
		panic(err)
	}

	err1 := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic1, Partition: kafka.PartitionAny},
		Value:          value1,
	}, nil)

	err2 := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic2, Partition: kafka.PartitionAny},
		Value:          value2,
	}, nil)

	if err1 != nil || err2 != nil {
		fmt.Println(err1.Error(), err2.Error())
		err = producer.AbortTransaction(ctx)
		if err != nil {
			panic(err)
		}
	}
	err = producer.CommitTransaction(ctx)
	if err != nil {
		fmt.Println("failed to commit transaction", err.Error())
		err = producer.AbortTransaction(ctx)
		if err != nil {
			panic(err)
		}
	}

	//Abort Transaction example
	err = producer.BeginTransaction()
	if err != nil {
		panic(err)
	}

	err1 = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic1, Partition: kafka.PartitionAny},
		Value:          value1,
	}, nil)

	err2 = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic2, Partition: kafka.PartitionAny},
		Value:          value2,
	}, nil)

	if err1 != nil || err2 != nil {
		fmt.Println(err1.Error(), err2.Error())
		err = producer.AbortTransaction(ctx)
		if err != nil {
			panic(err)
		}
	}
	err = producer.AbortTransaction(ctx)
	if err != nil {
		fmt.Println("failed to abort transaction", err.Error())
		panic(err)
	}

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
		producer.AbortTransaction(ctx)
		producer.Close()
	}()
}
