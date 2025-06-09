package producer

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

func GenerateRandomDateTime() time.Time {
	start := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	randSeconds := rand.Int63n(end.Unix() - start.Unix())
	return start.Add(time.Duration(randSeconds) * time.Second)
}

func GenerateRandomInteger(from int, to int, rnd rand.Rand) int {
	return from + rnd.Intn(to-from) + 1
}

// func GenerateRandomString(length int) string {
// 	var str string = ""
// 	const chars = "abcdef ghijklmnopqr stuvwxyz "
// 	for int(len(str)) < length {
// 		str = str + string(chars[GenerateRandomInteger(0, len(chars)-1)])
// 	}
// 	return str
// }

func GetRandomString(length int) map[string]interface{} {

	file, err := os.Open("dummy_content.json")
	if err != nil {
		log.Fatalf("Failed to open file : %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)

	var result map[string]interface{}
	if err := json.Unmarshal(byteValue, &result); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	return result
}
