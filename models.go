package main

type Quiz struct {
	Questions []Question
}

type Question struct {
	Text    string
	Answers []Answer
}

type Answer struct {
	Text    string
	Correct bool
}

type User struct {
	Uid   int
	Email string
}
