package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
1. Connect to a kafka topic
2. For each tweet field, generate a random value
3. In an infinite loop -> send payloads to the kafka topic
		- Figure out worker count
4. Produce logs to the logs/
*/

type Tweet struct {
	id         uint32
	text       string
	author_id  uint32
	retweets   uint32
	likes      uint32
	created_at time.Time // ?
	edited     bool
}

func main() {
	generated := generateRandomDateTime()
	fmt.Println(generated)
}

func generateRandomDateTime() time.Time {
	start := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	randSeconds := rand.Int63n(end.Unix() - start.Unix())
	return start.Add(time.Duration(randSeconds) * time.Second)
}

// Before writing the next function -> Respect TDD !
