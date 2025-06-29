package producer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ---------------------------------------- PRIVATE -------------------------------------------------------

func generateRandomDateTime() time.Time {
	start := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	randSeconds := rand.Int63n(end.Unix() - start.Unix())
	return start.Add(time.Duration(randSeconds) * time.Second)
}

func generateRandomInteger(from int, to int, rnd rand.Rand) int {
	return from + rnd.Intn(to-from) + 1
}

func getRandomString(key int) string {
	var allPossibleMessages = map[int]string{
		1:  "Be aware that system clock changes (e.g., time zone changes, NTP synchronization) can affect the accuracy of the measurements.",
		2:  "The methods above measure \"wall time\" (real-world elapsed time), which includes time spent waiting for I/O or other operations. If you need to measure CPU time specifically, you might need platform-specific methods.",
		3:  "When measuring concurrent functions (goroutines), you need to be careful about how you measure the time. The total time of the program might not be the sum of the individual goroutines' execution times.",
		4:  "For more accurate and detailed performance analysis, use Go's built-in benchmarking tools with the testing package.",
		5:  "Go maps are typed, so the value type must be consistent or use interface{} to support multiple types.",
		6:  "You can define a map with map[string]string or map[string]int if you know all values are of one type.",
		7:  "interface{} is Go's way to allow any type, similar to Python's dynamic typing.",
		8:  "We capitalize the first letter to export a specific function/variable and that makes it available to use for other packages ( makes it global ).",
		9:  "Goroutine - Lightweight thread managed by the Go runtime",
		10: "Threads are heavyweight and need more memory allocated to it.",
	}
	return allPossibleMessages[key]
}

// ----------------------------------------- PUBLIC ------------------------------------------------------

func CreateProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",

		/*new stuff*/
		"queue.buffering.max.messages": 10000,
		"queue.buffering.max.kbytes":   512000,
		"queue.buffering.max.ms":       5,
		"batch.num.messages":           1000,
		"message.send.max.retries":     3,
		"acks":                         1,
		"linger.ms":                    3,
		// "compression_type":				"gzip or snappy?",

		// "delivery.report.only.error": true, // Set this to reduce report volume -> This crashes!  Prints: Failed to create producer!
	})
	if err != nil {
		fmt.Println("Failed to create producer!")
	}
	return producer
}

func CreateMessage() map[string]any {
	var edited = [2]bool{true, false}
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	tweet := map[string]any{
		"id":         generateRandomInteger(1000000, 10000000, *seed),
		"text":       getRandomString(generateRandomInteger(1, 10+1, *seed)),
		"author_id":  generateRandomInteger(100000, 1000000, *seed),
		"likes":      generateRandomInteger(1000, 100000, *seed),
		"retweets":   generateRandomInteger(100, 5000, *seed),
		"created_at": generateRandomDateTime(),
		"edited":     edited[generateRandomInteger(0, 1, *seed)],
	}
	return tweet
}

func SendToKafkaTopic(message map[string]any, messageIndex int, topicName string, kafkaProducer *kafka.Producer) {
	formatted, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	err = kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: int32(-1)}, // equivalent to kafka.PartitionAny
		Value:          formatted,
	}, nil)

	if err != nil {
		fmt.Println("Error while sending message:", err)
	}

}
