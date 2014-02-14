package main

// {
//     "title": "Main Quiz",
//     "questions": [
//         {
//             "text": "What is the best programming language?",
//             "correct": 2,
//             "answers": [
//                 "Go",
//                 "Javascript",
//                 "Anything except .NET"
//             ],
//             "$$hashKey": "01L"
//         },
//         {
//             "text": "Is this the quiz you're looking for?",
//             "correct": 0,
//             "answers": [
//                 "This is not the quiz we are looking for",
//                 "Yes",
//                 "No"
//             ],
//             "$$hashKey": "01Q"
//         }
//     ]
// }

type Quiz struct {
	Title     string         `json:"title"`
	Questions []Question     `json:"questions"`
	Grades    map[string]int `json:"grades"` //map[sid]0-100
}

type Question struct {
	Text    string   `json:"text"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"` //offset in []Answers
}

type User struct {
	Uid   int
	Email string
}

type UserReply struct {
	sid string
	ans int
}

type Session struct {
	qid      int
	replies  chan UserReply
	state    chan int
	students map[string]string
}
