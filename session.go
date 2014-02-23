/*
session.go

writeErr
writeSuccess
handleAnswer
changeState
quizSesh
quit
*/


package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

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

// handleAnswer qets the quizID from the given URL w,
// gets an answer from a client, and stores it in a map.
//PUT/POST /quiz/{id}/answer
//  body (json):
//    {
//      "Id"      : string
//      "Answer"  : int
//    }
//
//  reply (json):
//    {
//      "Success": bool
//    }
//
//TODO auth student
func handleAnswer(w http.ResponseWriter, r *http.Request) {
	_ = sql.ErrNoRows
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	t := struct {
		Answer int    `json:"answer"`
		Id     string `json:"id"`
	}{}
	err := decoder.Decode(&t)
	if writeErr(err, w) {
		return
	}

	qid, err := strconv.Atoi(vars["id"])
	if writeErr(err, w) {
		return
	}

	if _, ok := qzSesh[qid]; !ok {
		writeErr(fmt.Errorf("Quiz session does not exist"), w)
		return
	}
	if _, ok := qzSesh[qid].students[t.Id]; !ok {
		writeErr(fmt.Errorf("User ID not in session"), w)
		return
	}
	qzSesh[qid].replies <- UserReply{t.Id, t.Answer}
	writeSuccess(w)
}

// must have cookie
//PUT /quiz/{id}/state
//  body:
//    state : int
//
// TODO json me, authenticate cookie
func changeState(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		return
	}

	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["id"])
	if err != nil {
		if writeErr(err, w) {
			return
		}
	}
	_, err = db.Exec(`SELECT qid 
  FROM classes, quiz 
  WHERE quiz.qid = $1 
  AND classes.qid = quiz.qid 
  AND classes.uid = $2`, qid, user.Uid)

	if writeErr(err, w) {
		return
	}
	decoder := json.NewDecoder(r.Body)
	t := struct {
		State int `json:"state"`
	}{}
	err = decoder.Decode(&t)
	if writeErr(err, w) {
		return
	}
	if _, ok := qzSesh[qid]; !ok && t.State == 0 {
		var results string
		err := db.QueryRow(`SELECT students 
                        FROM classes, quiz 
                        WHERE quiz.qid= $1 
                        AND classes.cid = quiz.cid`, qid).Scan(&results)
		var students map[string]string
		err = json.Unmarshal([]byte(results), &students)
		if writeErr(err, w) {
			return
		}

		qzSesh[qid] = Session{qid, make(chan UserReply), make(chan int), students}
		go quizSesh(qzSesh[qid])
	} else if !ok {
		writeErr(fmt.Errorf("Quiz session does not exist"), w)
		return
	}

	qzSesh[qid].state <- t.State
	writeSuccess(w)
}

type Response map[string]interface{}

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
	//[]map[sid]answer //TODO fix this size? we know how many ?'s there are...
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
			if state >= len(answers) { //only add to []answers as needed ~Dicey
				answers = append(answers, make(map[string]int))
			}
		}
	}
}

//map[sid]answer
func quit(qid int, qa []map[string]int) {
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

	//add number of correct
	// for each question
	//   for each answer
	for i, question := range qa { //TODO could insert qa into Quiz for further statistics
		for s, ans := range question {
			if ans == quiz.Questions[i].Correct {
				correct[s]++
			} else {
				if _, ok := correct[s]; !ok {
					correct[s] = 0
				}
			}
		}
	}

	//map[sid]0-100
	grades := make(map[string]int)
	for s, c := range correct {
		grades[s] = int(float64(c) / float64(len(quiz.Questions)) * 100)
	}

	quiz.Grades = grades
	q, err := json.Marshal(quiz)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(q))

	_, err = db.Exec(`UPDATE quiz SET info = $1 WHERE qid = $2`, string(q), qid)
	if err != nil {
		//TODO also really bad
		fmt.Println("cannot save grades", err)
	}

	//remove session
	delete(qzSesh, qid)
}
