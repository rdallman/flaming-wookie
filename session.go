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
func writeErr(err error, w http.ResponseWriter) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, Response{"success": false})
	}
}

func writeSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true})
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
		Answer int
		Id     string
	}{}
	err := decoder.Decode(&t)
	writeErr(err, w)

	qid, err := strconv.Atoi(vars["id"])
	writeErr(err, w)

	//TODO if session doesn't exist, reply with 401? something that indicates not in progress?
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
