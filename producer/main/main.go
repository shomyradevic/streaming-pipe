package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"producer/producer/producer"
	"time"
)

/*
1. Connect to a kafka topic
2. For each tweet field, generate a random value
3. In an infinite loop -> send payloads to the kafka topic
		- Figure out worker count
4. Produce logs to the logs/
*/

func sendToKafkaTopic(topicName string) {

	var edited = [2]bool{true, false}

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	tweet := map[string]any{
		"id":   producer.GenerateRandomInteger(1000000, 10000000, *seed),
		"text": producer.GetRandomString(producer.GenerateRandomInteger(1, 10+1, *seed)),
		// This needs to be fixed - We are doing I/O at every iteration, which doesnt make sense. It takes ~ 17 secs for 100K messages, which is a lot.

		"author_id":  producer.GenerateRandomInteger(100000, 1000000, *seed),
		"retweets":   producer.GenerateRandomInteger(10000, 100000, *seed),
		"likes":      producer.GenerateRandomInteger(100, 1000, *seed),
		"created_at": producer.GenerateRandomDateTime(),
		"edited":     edited[producer.GenerateRandomInteger(0, 1, *seed)]}

	formatted, err := json.MarshalIndent(tweet, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("tweet:", string(formatted))
}

func run(workerCount uint8, messageCount int) {
	fmt.Println("Using workers: ", workerCount)
	for i := 0; i < messageCount; i++ {
		sendToKafkaTopic("raw_topic")
	}
}

func main() {
	defer fmt.Println("Program finished with execution!") // :)

	start := time.Now()

	run(1, 100000)

	elapsed := time.Since(start)
	fmt.Println("Execution time: ", elapsed)
}
