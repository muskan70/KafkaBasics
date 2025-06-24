package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProductInventory struct {
	Id       int
	Name     string
	Quantity int
}

func producerWorker(producer *kafka.Producer, messages chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	topic := KafkaTopic
	for msg := range messages {
		SendMessageToKafka(msg, topic)
	}
}

func RunMultiThreadProducer() {
	const numWorkers = 5
	messages := make(chan []byte, 100) // Buffered channel for messages
	var wg sync.WaitGroup

	// Start producer workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go producerWorker(producer, messages, &wg)
	}

	// Produce messages
	for i := range 10000 {
		productName := "productName_" + strconv.Itoa(i)
		productInv := ProductInventory{
			Id:       i,
			Name:     productName,
			Quantity: i,
		}
		value, err := json.Marshal(productInv)
		if err != nil {
			panic(err)
		}
		messages <- value
	}
	close(messages) // Signal workers to finish

	wg.Wait()                 // Wait for all workers to complete
	producer.Flush(15 * 1000) // Flush any remaining messages with a timeout
	fmt.Println("Producer finished.")
}
