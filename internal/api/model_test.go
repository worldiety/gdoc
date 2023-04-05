package api

func ExampleList() {
	myList := List[string]{}
	myList.Add("hello world")
}

func ExampleList_Add() {
	myList := List[string]{}
	myList.Add("add example")
}
