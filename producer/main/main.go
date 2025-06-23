package main

import (
	"fmt"
	"producer/producer/producer"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func run(workerCount uint8, messageCount int, kafkaProducer *kafka.Producer, topicName string) {
	fmt.Println("Using workers: ", workerCount)
	for i := 0; i < messageCount; i++ {
		message := producer.CreateMessage()
		producer.SendToKafkaTopic(message, i, topicName, kafkaProducer)
	}
}

func main() {
	var workerCount uint8 = 1
	var messageCount int = 100000
	var kafkaProducer = producer.CreateProducer()
	var topicName string = "raw-topic"

	defer kafkaProducer.Flush(5000)
	defer kafkaProducer.Close()

	start := time.Now()

	run(workerCount, messageCount, kafkaProducer, topicName)

	elapsed := time.Since(start)
	fmt.Println("Execution time: ", elapsed)
}
