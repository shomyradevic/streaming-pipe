package main

import (
	"fmt"
	"producer/producer"
)

func main() {
	generated := producer.GenerateRandomDateTime()
	fmt.Println(generated)
}
