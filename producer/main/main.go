package main

import (
	"fmt"
	"os"
	"os/signal"
	"producer/producer/producer"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func run(workerCount uint8, messageCount int, kafkaProducer *kafka.Producer, topicName string) {
	fmt.Println("Using workers: ", workerCount)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	/*
		Looks like this goroutine doesnt prevent producer from just silently stopping at about 30K messages...
		Maybe its a different signal? We need to investigate whats actually happening and why its just casually stopping...
	*/
	go func() {
		<-sigchan
		fmt.Println("Flushing...")
		kafkaProducer.Flush(5000)
		kafkaProducer.Close()
		os.Exit(0)
	}()

	for i := 0; i < messageCount; i++ {
		message := producer.CreateMessage()
		producer.SendToKafkaTopic(message, i, topicName, kafkaProducer)
	}
	kafkaProducer.Flush(5000)
	kafkaProducer.Close()
}

func main() {
	var workerCount uint8 = 1
	var messageCount int = 100000
	var kafkaProducer = producer.CreateProducer()
	var topicName string = "raw-topic"

	// defer kafkaProducer.Flush(5000)
	// defer kafkaProducer.Close()

	start := time.Now()

	run(workerCount, messageCount, kafkaProducer, topicName)

	elapsed := time.Since(start)
	fmt.Println("Execution time: ", elapsed)
}
