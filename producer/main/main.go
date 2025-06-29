package main

import (
	"fmt"
	"os"
	"os/signal"
	"producer/producer/producer"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func run(workerCount uint8, messageCount int, kafkaProducer *kafka.Producer, topicName string) {

	// Trap SIGINT/SIGTERM so we can flush before exiting
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	var wg sync.WaitGroup // This can be important

	fmt.Println("Using workers: ", workerCount)

	go func() {
		for e := range kafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < messageCount; i++ {
			message := producer.CreateMessage()
			producer.SendToKafkaTopic(message, i, topicName, kafkaProducer)
		}
	}()
	wg.Wait()

	fmt.Printf("Finished with loop!")
	remaining := kafkaProducer.Flush(5000)
	fmt.Printf("Remaining after flush: %d", remaining)

	// fmt.Println("🏁 All done. Sleeping 2s before exit to flush OS/Docker buffers...")
	/* With this, everything is printed when using 5000 messages, which means that docker is early exiting.
	However, when doing 200K messages, it doesnt print everything, because it early exits.*/
	// time.Sleep(5 * time.Second)

	kafkaProducer.Close()

}

func main() {
	var workerCount uint8 = 1
	var messageCount int = 1000000
	var kafkaProducer = producer.CreateProducer()
	var topicName string = "raw-topic"

	start := time.Now()

	run(workerCount, messageCount, kafkaProducer, topicName)

	elapsed := time.Since(start)
	fmt.Println("Execution time: ", elapsed)

}
