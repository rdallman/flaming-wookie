package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: need to parse body for the title of the quiz, for now everything is just Test
// TODO: unsafe, definitely don't need to print Postgres errors back to the browser, but need debugging
// TODO: flash message to show quiz was added, and redirect
//
//QUIZ info of the form:
//  info :
//    questions : [
//      {
//        text : string,
//        answers : [
//          string,
//          ...
//        ],
//        correct : string
//      },
//      ...
//    ],
//    grades : {
//      studentid (int) : grade (int)
//    }
//
// on creation just make a blank map for grades

// inserts a quiz from an AJAX POST request
func handleQuizCreate(w http.ResponseWriter, r *http.Request) {
	// grab body of request (should be the json of the quiz)
	var quizJson Quiz
	p := make([]byte, r.ContentLength)
	_, err := r.Body.Read(p)
	if err != nil {
		// return an *actual* error
		fmt.Fprintf(w, "I'm illiterate...")
	} else {
		err = json.Unmarshal(p, &quizJson)

		// insert the quiz
		err = insertQuiz(quizJson.Title, p)
		if err != nil {
			fmt.Fprintf(w, err.Error()) // remove eventually, needed for debugging
		}
	}

}

func insertQuiz(title string, quizData []byte) error {
	_, err := db.Exec(`INSERT INTO quiz (title, info, cid) 
		VALUES($1, $2, $3)`, title, quizData, 1)
	if err != nil {
		return err
	}
	return nil

}
