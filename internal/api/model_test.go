package api

import "fmt"

func ExampleList() {
	myList := List[string, string, Constant]{}
	myList.Add("hello world")
}

func ExampleList_Add() {
	myList := List[string, int, Constant]{}
	myList.Add("add example")
	fmt.Print("Hello")
	// Output: "Hello"
}
