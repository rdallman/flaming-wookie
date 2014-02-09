package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

// handleAnswer qets the quizID from the given URL w,
// gets an answer from a client, and stores it in a map.
//PUT/POST /quiz/{id}/answer
//  body:
//    answer : string
//
//TODO auth student
func handleAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(r.Body)
	qid, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//TODO Basic Auth is base64 encoded and gross... when you're feeling extra bored
	sid := r.Header.Get("Authorization")
	sid = sid[:len(sid)-1] //slice off last
	fmt.Fprintf(w, "%d %v", sid, err)
	a, err := strconv.Atoi(r.FormValue("answer"))
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
	answers := make([]map[string]int, 0)
	for {
		select {
		case ur := <-s.replies:
			fmt.Println(ur.sid, ur.ans)
			answers[state][ur.sid] = ur.ans
		case state = <-s.state:
			fmt.Println(state)
			if state < 0 {
				go quit(s.qid, answers)
				break
			}
			if state >= len(answers) {
				answers = append(answers, make(map[string]int))
			}
		}
	}
}

//map[sid]answer
func quit(qid int, qa []map[string]int) {
	//var quiz Quiz
	var qstring string
	err := db.QueryRow(`SELECT info FROM quiz WHERE qid = $1`, qid).Scan(&qstring)
	if err != nil {
		//TODO uh this is really bad at this point
		fmt.Println("cannot find quiz", err, qid)
	}

	var quiz Quiz
	err = json.Unmarshal([]byte(qstring), &quiz)
	if err != nil {
		fmt.Println(err)
	}

	//map[sid]#correct
	correct := make(map[string]int)

	for i, q := range qa { //TODO could insert qa into Quiz for further statistics
		fmt.Println(i)
		for s, ans := range q {
			fmt.Println(s)
			if ans == quiz.Questions[i].Correct {
				correct[s]++
			} else {
				if _, ok := correct[s]; !ok {
					correct[s] = 0
				}
			}
		}
	}

	//map[sid]%
	grades := make(map[string]int)
	for s, c := range correct {
		grades[s] = c / len(quiz.Questions)
	}

	quiz.Grades = grades
	q, err := json.Marshal(quiz)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(q, string(q))

	_, err = db.Exec(`UPDATE quiz SET info = $1 WHERE qid = $2`, string(q), qid)
	if err != nil {
		//TODO also really bad
		fmt.Println("cannot save grades", err)
	}

	delete(qzSesh, qid)
}
