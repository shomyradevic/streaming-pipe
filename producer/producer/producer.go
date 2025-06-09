package producer

import (
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

// type Tweet struct {
// 	id         uint32
// 	text       string
// 	author_id  uint32
// 	retweets   uint32
// 	likes      uint32
// 	created_at time.Time // ?
// 	edited     bool
// }

func GenerateRandomDateTime() time.Time {
	start := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	randSeconds := rand.Int63n(end.Unix() - start.Unix())
	return start.Add(time.Duration(randSeconds) * time.Second)
}

func GenerateRandomInteger(from int, to int) int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return from + rnd.Intn(to-from) + 1
}

func GenerateRandomString(length int) string {
	var str string = ""
	const chars = "abcdef ghijklmnopqr stuvwxyz "
	for int(len(str)) < length {
		str = str + string(chars[GenerateRandomInteger(0, len(chars)-1)])
	}
	return str
}

// Before writing the next function -> Respect TDD !
