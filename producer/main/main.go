package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"producer/producer/producer"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func createMessage() map[string]any {
	var edited = [2]bool{true, false}
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	tweet := map[string]any{
		"id":         producer.GenerateRandomInteger(1000000, 10000000, *seed),
		"text":       producer.GetRandomString(producer.GenerateRandomInteger(1, 10+1, *seed)),
		"author_id":  producer.GenerateRandomInteger(100000, 1000000, *seed),
		"retweets":   producer.GenerateRandomInteger(10000, 100000, *seed),
		"likes":      producer.GenerateRandomInteger(100, 1000, *seed),
		"created_at": producer.GenerateRandomDateTime(),
		"edited":     edited[producer.GenerateRandomInteger(0, 1, *seed)]}
	return tweet
}

func sendToKafkaTopic(message map[string]any, messageIndex int, topicName string, kafkaProducer *kafka.Producer) {
	formatted, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	err = kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: int32(-1)}, // equivalent to kafka.PartitionAny
		Value:          formatted,
	}, nil)

	if err == nil {
		fmt.Println("Message sent:", messageIndex)
	} else {
		fmt.Println("Error while sending %d. message:", err)
	}
}

func run(workerCount uint8, messageCount int, kafkaProducer *kafka.Producer, topicName string) {
	fmt.Println("Using workers: ", workerCount)
	for i := 0; i < messageCount; i++ {
		message := createMessage()
		sendToKafkaTopic(message, i, topicName, kafkaProducer)
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
