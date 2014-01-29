package main

type Quiz struct {
	Questions []Question  `json:questions`
	Grades    map[int]int `json:grades` //map[sid]0-100
}

type Question struct {
	Text    string
	Answers []Answer
	Correct string
}

type Answer struct {
	Text string
}

type User struct {
	Uid   int
	Email string
}

type QAnswers struct {
	//map[sid]answer#
	studentAnswer map[int]string
}

type UserReply struct {
	sid int
	ans string
}

type Session struct {
	qid     int
	replies chan UserReply
	state   chan int
}
