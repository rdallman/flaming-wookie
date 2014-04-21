package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	TRACE   *log.Logger //standard stuff
	INFO    *log.Logger //special information
	WARNING *log.Logger //there is something you need to know about
	ERROR   *log.Logger //something has failed
)

func init() {
	//create log file
	file, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}

	TRACE = log.New(io.MultiWriter(file, ioutil.Discard),
		"TRACE: ",
		log.LstdFlags)

	INFO = log.New(io.MultiWriter(file, os.Stdout),
		"INFO: ",
		log.LstdFlags)

	WARNING = log.New(io.MultiWriter(file, os.Stdout),
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ERROR = log.New(io.MultiWriter(file, os.Stdout),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

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
	students map[string]bool
}

//checks for err, replies with false success reply
func writeErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return true
	}
	return false
}

//TODO want support for multiple items?
func writeSuccess(w http.ResponseWriter, info ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	r := Response{"success": true}
	if len(info) == 1 {
		r["info"] = info[0]
	}
	fmt.Fprint(w, r)
}

type Response map[string]interface{}

func (r Response) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}
