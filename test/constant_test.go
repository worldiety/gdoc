package test

import (
	"fmt"
)

// A is a test Constant and this is a test comment.
const A = "20"

const B = 10
const (
	Test = "ABC"
)

var ReadError error

type TestType string

// A test variable
var S TestType = "hello"

// testFunction
func FuncTest(testP string) {
	testP = "world"
	fmt.Sprint(string(S) + testP)
}
