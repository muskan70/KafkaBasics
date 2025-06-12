package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	KafkaServer = "localhost:9092,localhost:9093,localhost:9094"
	KafkaTopic  = "stock-topic"
)

type Stock struct {
	Id        int
	Name      string
	LastPrice float64
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
		"client.id":         "myProducer",
		"acks":              "all",
	})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	topic := KafkaTopic
	delivery_chan := make(chan kafka.Event, 10000)

	for i := range 10000 {
		stockName := "stockName_" + strconv.Itoa(i)
		stock := Stock{
			Id:        i,
			Name:      stockName,
			LastPrice: 1000.00 + float64(i),
		}
		value, err := json.Marshal(stock)
		if err != nil {
			panic(err)
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
		}, delivery_chan)

		if err != nil {
			panic(err)
		}

		e := <-delivery_chan
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}

	}
	close(delivery_chan)
}
