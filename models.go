package main

type Quiz struct {
	Title     string      `json:title`
	Questions []Question  `json:questions`
	Grades    map[int]int `json:grades` //map[sid]0-100
}

type Question struct {
	Text    string   `json:text`
	Answers []string `json:answers`
	Correct int      `json:correct` //offset in []Answers
}

type User struct {
	Uid   int
	Email string
}

type UserReply struct {
	sid int
	ans int
}

type Session struct {
	qid     int
	replies chan UserReply
	state   chan int
}
