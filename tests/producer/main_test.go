package producer_test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGenerateRandomDateTime(t *testing.T) {
	/*
		This test needs to check is the function:
		- returning a time.Time type
		- Also need to parameterize this better in order to simplify adding test cases. Look for something cutting-edge.
	*/
	result := ""
	varType := reflect.TypeOf(result).Kind()
	fmt.Println("vartype is ", varType)
}
