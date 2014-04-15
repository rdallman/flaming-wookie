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
	state      int                      // current question or INFINITY
	questions  []sendQ                  // ready to send out questions
	registered map[string]bool          // students who can connect
	conns      map[*websocket.Conn]bool // students connected
	answers    []map[string]int         // []map[sid]answer where i = ?#
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

		sesh.answers[sesh.state][sid] = data.Answer
	}
}

// validate that a teacher owns requested quiz, spin up a session
func teachServer(ws *websocket.Conn) {
	user := auth(ws.Request())
	if user == nil {
		return
	}

	qid := mux.Vars(ws.Request())["id"]

	var count int
	err := db.QueryRow(`
  SELECT COUNT(classes.uid)
  FROM classes, quiz
  WHERE quiz.qid = $1
  AND classes.uid = $2`, qid, user.Uid).Scan(&count)

	if err != nil || count < 1 { // not authorized or something terribly wrong
		return
	}

	//if _, ok := sessions[qid]; !ok { // TODO
	cid := newSession(qid)
	listenTeach(ws, cid)
}

// if there's currently a session for this quiz, overwrites
// mostly pulling data from db to persist in a *session
func newSession(qid string) string {
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
		questions:  qs,
		state:      math.MaxInt32,
		registered: students,
		conns:      make(map[*websocket.Conn]bool),
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
	var qstring string
	sesh := sessions[cid]
	answers := sesh.answers
	err := db.QueryRow(`SELECT info FROM quiz WHERE qid = $1`, sesh.qid).Scan(&qstring)
	if err != nil {
		//TODO uh this is really bad at this point
		fmt.Println("cannot find quiz", err, sesh.qid)
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
}
