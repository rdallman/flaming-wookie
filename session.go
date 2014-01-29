package main

import (
	"encoding/json"
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
	}
	//TODO Basic Auth is base64 encoded and gross... when you're feeling extra bored
	sid, err := strconv.Atoi(r.Header.Get("Authorization"))
	fmt.Fprintf(w, "%d %v", sid, err)
	a := r.FormValue("answer")
	//TODO if session doesn't exist, reply with 401? something that indicates not in progress?
	qzSesh[qid].replies <- UserReply{sid, a}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{true})
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
	fmt.Println(state)
	qzSesh[qid].state <- state
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{true})
	//TODO reply w/ 200
}

type Response struct {
	Success bool `json:success`
}

func (r Response) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

//send 0
func quizSesh(s Session) {
	state := 0
	//[]map[sid]answer
	answers := make([]map[int]string, 0)
	for {
		select {
		case ur := <-s.replies:
			fmt.Println(ur.sid, ur.ans)
			if state >= len(answers) {
				answers = append(answers, make(map[int]string))
			}
			answers[state][ur.sid] = ur.ans
		case state = <-s.state:
			fmt.Println(state)
			if state < 0 {
				go quit(s.qid, answers)
				break
			}
		}
	}
}

//map[sid]answer
func quit(qid int, qa []map[int]string) {
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
		for s, ans := range q {
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
