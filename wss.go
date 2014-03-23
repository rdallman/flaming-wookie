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
	state      int
	questions  []sendQ
	registered map[string]bool
	answers    []map[string]int
	conns      map[*websocket.Conn]string
}

type students map[*websocket.Conn]bool

type sendQ struct {
	Text    string   `json:"text"`
	Answers []string `json:"answers"`
}

func init() {
	sessions = make(map[string]*session)
}

// Echo the data received on the WebSocket.
func studServer(ws *websocket.Conn) {
	id := mux.Vars(ws.Request())["id"]
	s, ok := sessions[id]

	if !ok {
		ws.Close()
		return
	}

	fmt.Println(ws.Request().Header["Authorization"])

	return

	sid := ws.Request().URL.User.Username()
	if _, ok := s.registered[sid]; !ok {
		ws.Close()
		return
	}

	if s.state == math.MaxInt64 {
		websocket.JSON.Send(ws, &sendQ{Text: `{"text":"quiz not started"}`, Answers: nil})
	} else {
		websocket.JSON.Send(ws, s.questions[s.state])
	}
}

func listenAnswers(ws *websocket.Conn, sid, cid string) {
	for {
		var data struct {
			Answer int `json:"answer"`
		}
		err := websocket.JSON.Receive(ws, &data)
		if err != nil { // io.EOF = disconnect TODO separately
			ws.Close()
			return
		}
		sesh := sessions[cid]

		sesh.answers[sesh.state][sid] = data.Answer

		websocket.JSON.Send(ws, &data)
	}
}

func teachServer(ws *websocket.Conn) {
	// TODO authorize teachers

	qid := mux.Vars(ws.Request())["id"]
	// TODO yeahh they give QID and sessions are CID based
	// so find another way to check while authorizing
	//
	//if _, ok := sessions[qid]; !ok {
	cid := newSession(qid)
	listenTeach(ws, cid)
	//}
}

// returns cid
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

	sessions[cid] = &session{
		qid:        qid,
		title:      title,
		questions:  qs,
		state:      math.MaxInt64,
		registered: students,
		conns:      make(map[*websocket.Conn]string),
	}

	return cid
}

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

		if schange.State < 0 {
			// TODO quit, close all conns and save grades
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
				fmt.Println(err)
			}
		}
	}
}
