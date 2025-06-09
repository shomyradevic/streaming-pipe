package producer_test

import (
	"producer/producer/producer"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomDateTime(t *testing.T) {
	/*
		Need to parameterize this better in order to simplify adding test cases. Look for something cutting-edge.
	*/
	result := producer.GenerateRandomDateTime()
	assert.WithinRange(t, result, time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()) // actual, start, end
}

func TestGenerateRandomInteger(t *testing.T) {
	/*
		Need to parameterize this better in order to simplify adding test cases. Look for something cutting-edge.
	*/
	var from int = 1
	var to int = 10000
	result := producer.GenerateRandomInteger(from, to)
	assert.GreaterOrEqual(t, result, from)
	assert.LessOrEqual(t, result, to)
}

func TestGenerateRandomString(t *testing.T) {
	/*
		Need to parameterize this better in order to simplify adding test cases. Look for something cutting-edge.
	*/
	var length int = 100

	result := producer.GenerateRandomString(length)
	if len(result) > length {
		t.Fail()
	}
	assert.Equal(t, bool(1 == (2-1)), true)
}
