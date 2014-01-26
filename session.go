package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

//PUT/POST /quiz/{id}/answer
//  body:
//    answer : string
//
//TODO auth student
func handleAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//TODO does this fallthrough?
	}
	sid, err := strconv.Atoi(r.Header.Get("Authorization")) //TODO yeah this has a :
	fmt.Fprintf(w, "%d %v", sid, err)
	a := r.FormValue("answer")
	//TODO if session doesn't exist, reply with 401? something that indicates not in progress?
	qzSesh[qid].replies <- UserReply{sid, a}
}

// must have cookie
//PUT /quiz/{id}/state
//  body:
//    state : int
func changeState(w http.ResponseWriter, r *http.Request) {
	//auth()
	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//TODO does this fallthrough?
	}
	state, err := strconv.Atoi(r.FormValue("state"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if _, ok := qzSesh[qid]; !ok && state == 0 {
		qzSesh[qid] = Session{qid, make(chan UserReply), make(chan int)}
		go quizSesh(qzSesh[qid])
	}
	qzSesh[qid].state <- state
	//TODO reply w/ 200
}

//send 0
func quizSesh(s Session) {
	var state int
	s.state <- 0
	answers := make([]QAnswers, 0)
	for {
		select {
		case ur := <-s.replies:
			fmt.Println(ur.sid, ur.ans)
			answers[state].studentAnswer[ur.sid] = ur.ans
		case state = <-s.state:
			fmt.Println(state)
			if state < 0 {
				go quit(s.qid, answers)
				break
			}
		}
	}
}

func quit(qid int, qa []QAnswers) {
	var quiz Quiz
	err := db.QueryRow(`SELECT info FROM quiz WHERE qid = $1`, qid).Scan(&quiz)
	if err != nil {
		//TODO uh this is really bad at this point
		fmt.Println("cannot find quiz")
	}

	//map[sid]#correct
	correct := make(map[int]int)

	for i, q := range qa {
		//TODO could insert qa into Quiz here for further statistics
		for s, ans := range q.studentAnswer {
			if ans == quiz.Questions[i].Correct {
				correct[s]++
			}
		}
	}

	grades := make(map[int]int)
	for s, c := range correct {
		grades[s] = c / len(quiz.Questions)
	}

	quiz.Grades = grades
	_, err = db.Exec(`UPDATE quiz SET info = $1 WHERE qid = $2`, quiz, qid)
	if err != nil {
		//TODO also really bad
		fmt.Println("cannot save grades")
	}

	delete(qzSesh, qid)
}
