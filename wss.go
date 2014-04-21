package main

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/gorilla/mux"

	"code.google.com/p/go.net/websocket"
)

// classid to session
var sessions map[string]*session

type session struct {
	qid        string
	title      string
	seshType   int
	state      int                      // current question or INFINITY
	questions  []sendQ                  // ready to send out questions
	registered map[string]bool          // students who can connect
	conns      map[*websocket.Conn]bool // students connected
	answers    []map[string]int         // []map[sid]answer where i = ?# (group by question)
	answersByStud map[string][]int 		// map[sid][]answer where i = ?# (group by student)
}

type sendQ struct {
	Text    string   `json:"text"`
	Answers []string `json:"answers"`
}

func init() {
	// map[cid]session
	sessions = make(map[string]*session)
}

// on new student connection, assert they are in class and then listen for their answers
// if there are any quizzes going on.
func studServer(ws *websocket.Conn) {
	id := mux.Vars(ws.Request())["id"]
	s, ok := sessions[id]

	if !ok {
		return
	}

	// TODO don't do this
	auth := ws.Request().Header["Authorization"]
	if len(auth) < 1 {
		return
	}
	sid := auth[0]

	// TODO put me back
	fmt.Println(sid, "connected to", id)
	if _, ok := s.registered[sid]; !ok {
		ws.Close()
		return
	}
	s.conns[ws] = true

	if s.state != math.MaxInt32 {
		websocket.JSON.Send(ws, s.questions[s.state])
	}

	listenAnswers(ws, sid, id)
}

// for ever take in answers from the student until teacher closes quiz
// TODO heartbeat
func listenAnswers(ws *websocket.Conn, sid, cid string) {
	for {
		var data struct {
			Answer int `json:"answer"`
		}
		err := websocket.JSON.Receive(ws, &data)
		if err != nil { // io.EOF = disconnect TODO separately
			delete(sessions[cid].conns, ws) // remove from session
			return
		}
		fmt.Println(sid, "says:", data.Answer)
		sesh := sessions[cid]

		// group by question
		sesh.answers[sesh.state][sid] = data.Answer
		
		// group by student
		if _, found := sesh.answersByStud[sid]; !found {
			sesh.answersByStud[sid] = make([]int, len(sesh.questions))
		}
		sesh.answersByStud[sid][sesh.state] = data.Answer
	}
}

// validate that a teacher owns requested quiz, spin up a session
func teachServer(ws *websocket.Conn) {
	user := auth(ws.Request())
	if user == nil {

		return
	}

	qid := mux.Vars(ws.Request())["id"]

	var qtype int
	err := db.QueryRow(`
  SELECT quiz.type
  FROM classes, quiz
  WHERE quiz.qid = $1
  AND classes.uid = $2`, qid, user.Uid).Scan(&qtype)

	if err != nil || qtype < 1 { // not authorized or something terribly wrong
		return
	}

	//if _, ok := sessions[qid]; !ok { // TODO
	cid := newSession(qid, qtype)
	listenTeach(ws, cid)
}

// if there's currently a session for this quiz, overwrites
// mostly pulling data from db to persist in a *session
func newSession(qid string, qtype int) string {
	//accepting == true
	var studjson, cid, title, qjson string
	err := db.QueryRow(`
      SELECT students, quiz.cid, info->>'title', info->>'questions'
      FROM classes, quiz
      WHERE quiz.qid = $1
      AND classes.cid = quiz.cid`, qid).Scan(&studjson, &cid, &title, &qjson)

	if err != nil {
		fmt.Println(err)
		return "" // bad
	}

	var sids []struct {
		Sid string `json:"sid"`
	}
	err = json.Unmarshal([]byte(studjson), &sids)
	if err != nil {
		fmt.Println(studjson)
		fmt.Println(err)
		return ""
	}
	students := make(map[string]bool)
	for _, s := range sids {
		students[s.Sid] = true
	}

	var qs []sendQ
	err = json.Unmarshal([]byte(qjson), &qs)

	// TODO(reed): yeah probably we don't do this here
	sessions[cid] = &session{
		qid:        qid,
		title:      title,
		seshType:	1,
		questions:  qs,
		state:      math.MaxInt32,
		registered: students,
		conns:      make(map[*websocket.Conn]bool),
		answersByStud: make(map[string][]int),
	}

	return cid
}

// listen for state changes from a teacher, close session and save grades when done
func listenTeach(ws *websocket.Conn, cid string) {
	for {
		var schange struct {
			State int `json:"state"`
		}
		err := websocket.JSON.Receive(ws, &schange)
		if err != nil { // io.EOF = disconnect
			fmt.Println(err)
			break
		} // TODO handle both io.EOF and json err
		sesh := sessions[cid]
		fmt.Println("class", cid, "in state:", schange.State)

		if schange.State < 0 {
			// TODO quit, close all conns and save grades
			for conn, _ := range sesh.conns {
				conn.Close()
			}
			go func() {
				gradeSesh(cid)
				delete(sessions, cid)
				// TODO(reed) heartbeat the delete
			}()
			return
		}

		if sesh.state != schange.State {
			sesh.state = schange.State
			// TODO check if in range
			qjson := sesh.questions[sesh.state]
			if sesh.state >= len(sesh.answers) {
				sesh.answers = append(sesh.answers, make(map[string]int))
			}
			for conn, _ := range sesh.conns {
				err := websocket.JSON.Send(conn, qjson)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

// grades a session and inserts results into db, 0-100 for each student
func gradeSesh(cid string) {
	sesh := sessions[cid]

	if sesh.seshType == 2 {
		// attendance
		attendance, err := json.Marshal(sesh.answers)
		if err != nil {
			fmt.Println("cannot marshal answers by question", err, sesh.qid)
		}
		_, err = db.Exec(`insert into attendance (cid, students, date_created) values($1, $2, now())`, cid, string(attendance))
		if err != nil {
			fmt.Println("cannot insert", err, sesh.qid)
		}
		return 

	}
	var qstring string
	var qtype int
	
	answers := sesh.answers
	err := db.QueryRow(`SELECT info, type FROM quiz WHERE qid = $1`, sesh.qid).Scan(&qstring, &qtype)
	if err != nil {
		//TODO uh this is really bad at this point
		fmt.Println("cannot find quiz", err, sesh.qid)
	}

	// dump into session_dump
	answersDump, err := json.Marshal(answers)
	if err != nil {
		fmt.Println("cannot marshal answers by question", err, sesh.qid)
	}

	answersByStudDump, err := json.Marshal(sesh.answersByStud)
	if err != nil {
		fmt.Println("cannot marshal answers by student", err, sesh.qid)
	}
	
	_, err = db.Exec(`insert into session_dump (qid, by_question, by_student, date_created) values($1, $2, $3, now())`, sesh.qid, string(answersDump), string(answersByStudDump))
	if err != nil {
		fmt.Println("cannot insert into session_dump", err, sesh.qid)
	}


	var quiz Quiz
	err = json.Unmarshal([]byte(qstring), &quiz)
	if err != nil {
		fmt.Println(err)
	}

	//map[sid]#correct
	correct := make(map[string]int)

	if qtype == 1 { // quiz
		//add number of correct
		// for each question
		//   for each answer
		for i, question := range answers { //TODO could insert qa into Quiz for further statistics
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

		_, err = db.Exec(`UPDATE quiz SET info = $1 WHERE qid = $2`, string(q), sesh.qid)
		if err != nil {
			//TODO also really bad
			fmt.Println("cannot save grades", err)
		}

	} else { // poll
		// this is just to put something in the grades key for polls
		grades := make(map[string]int)
		grades["filler"] = 0
		quiz.Grades = grades
		q, err := json.Marshal(quiz)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(q))

		_, err = db.Exec(`UPDATE quiz SET info = $1 WHERE qid = $2`, string(q), sesh.qid)
		if err != nil {
			//TODO also really bad
			fmt.Println("cannot save grades", err)
		}
	}

	

}


func attendanceServer(ws *websocket.Conn) {
	user := auth(ws.Request())
	if user == nil {
		return
	}

	cid := mux.Vars(ws.Request())["cid"]

	// can't use newSession here because we don't have qid, but we'll do the same stuff...
	var studentJson string
	err := db.QueryRow(`
      SELECT students
      FROM classes
      WHERE cid = $1`, cid).Scan(&studentJson)

	if err != nil {
		fmt.Println(err)
		return
	}

	var sids []struct {
		Sid string `json:"sid"`
	}
	err = json.Unmarshal([]byte(studentJson), &sids)
	if err != nil {
		fmt.Println(studentJson)
		fmt.Println(err)
		return
	}
	students := make(map[string]bool)
	for _, s := range sids {
		students[s.Sid] = true
	}

	attendanceQuestions := sendQ{Text: "Attendance", Answers: []string{"I'm here"}}

	sessions[cid] = &session{
		qid:        "1",	// just whatever, doesn't matter
		title:      "Attendance",
		seshType:	2,					// attendance
		questions:  []sendQ{attendanceQuestions},
		state:      math.MaxInt32,
		registered: students,
		conns:      make(map[*websocket.Conn]bool),
		answersByStud: make(map[string][]int),
	}
	listenTeach(ws, cid)
}