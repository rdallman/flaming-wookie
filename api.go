package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Expecting http body with JSON of form:
// {
//    name : string
//    students : {
//       id : name,
//       ...
//      }
//  }
//
func handleCreateClass(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		return //FIXME return error
	}

	j := struct {
		Name     string            `json:"name"`
		Students map[string]string `json:"students"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&j)
	if writeErr(err, w) {
		return
	}

	st, err := json.Marshal(j.Students)
	if writeErr(err, w) {
		return
	}

	// insert the quiz
	_, err = db.Exec(`INSERT INTO classes (name, uid, students) 
		VALUES($1, $2, $3)`, j.Name, user.Uid, string(st))
	if writeErr(err, w) {
		return
	}
	writeSuccess(w)
}

//add a student to a class
//
//
func handleAddStudents(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		return
	}

	cid := mux.Vars(r)["cid"]

	decoder := json.NewDecoder(r.Body)
	j := struct {
		name string
		sid  string
	}{}
	err := decoder.Decode(&j)
	if writeErr(err, w) {
		return
	}

	//TODO ->> json
	var jstring string
	err = db.QueryRow(`SELECT students FROM classes WHERE uid = $1 AND cid = $2`, user.Uid, cid).Scan(&jstring)
	if writeErr(err, w) {
		return
	}

	var students map[string]string
	err = json.Unmarshal([]byte(jstring), &students)
	if writeErr(err, w) {
		return
	}
	students[j.sid] = j.name

	js, err := json.Marshal(students)
	if writeErr(err, w) {
		return
	}

	_, err = db.Exec(`UPDATE classes SET students=$1 WHERE cid= $2`, string(js), cid)
	if writeErr(err, w) {
		return
	}

	writeSuccess(w)
}

// TODO: flash message to show quiz was added, and redirect
//
// Expecting http body with JSON of form:
//    {
//      questions : [
//        {
//          text : string,
//          answers : [
//            string,
//            ...
//          ],
//          correct : string
//        },
//        ...
//      ],
//      grades : {
//        studentid (int) : grade (int),
//        ...
//      }
//    }
//
// on creation just make a blank map for grades

// TODO add cid to frontend

// handleQuizCreate creates a quiz from an AJAX POST request
func handleQuizCreate(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		return
	}

	cid := mux.Vars(r)["cid"]

	var q Quiz
	err := json.NewDecoder(r.Body).Decode(&q)
	if writeErr(err, w) {
		return
	}

	info, err := json.Marshal(q)
	if writeErr(err, w) {
		return
	}

	// insert the quiz
	_, err = db.Exec(`
    INSERT INTO quiz (info, cid)
		VALUES($1, $2)
    `, string(info), cid) //TODO get CID from URL
	if writeErr(err, w) {
		return
	}
	writeSuccess(w)
}

// handleQuizGet qets the quizID from the end of the given URL w
// and writes the info json from the db back using r.
// ((add more later))
func handleQuizGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qID, err := strconv.Atoi(vars["id"])
	if writeErr(err, w) {
		return
	}
	//fmt.Printf("%d", qID) //testing
	var info string
	err = db.QueryRow(`SELECT info FROM quiz WHERE qid=$1`, qID).Scan(&info)
	if writeErr(err, w) {
		return
	}
	fmt.Fprintf(w, info)
}

// handleQuizUpdate qets the quizID from the end of the given URL w,
// gets the form data, and updates the quiz in the db.
// ((just an idea, not sure if we actually need this))

//TODO JSON me cap'n
func handleQuizUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//not sure we need all these, but for now...
	qid, err1 := strconv.Atoi(vars["id"])
	title := r.FormValue("title")
	info := r.FormValue("info")
	cid, err2 := strconv.Atoi(r.FormValue("cid"))

	if err1 != nil || err2 != nil {
		fmt.Println("err1 $1\terr2 $2", err1, err2)
	} else {
		_, err := db.Exec(`UPDATE quiz SET title=$1, info=$2, cid=$3 WHERE qid=$4`, title, info, cid, qid)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("\nUpdated.")
		}
	}
}

//Delete a quiz
func handleQuizDelete(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		return
	}
	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["id"])
	if writeErr(err, w) {
		return
	}
	_, err = db.Exec(`DELETE FROM quiz WHERE qid=$1`, qid)
	if writeErr(err, w) {
		return
	}
	writeSuccess(w)
}

func handleQuizList(w http.ResponseWriter, r *http.Request) {
	//title and id, return JSON
	auth := auth(r)
	if auth == nil {
		return
	}

	rows, err := db.Query(`
    SELECT qid, info->>'title'
    FROM quiz, classes
    WHERE classes.uid = $1
    AND classes.cid = quiz.cid
  `, auth.Uid)
	defer rows.Close()
	if writeErr(err, w) {
		return
	}
	quizzes := make(map[string]int)
	for rows.Next() {
		var qid int
		var title string
		err = rows.Scan(&qid, &title)
		if writeErr(err, w) {
			return
		}
		quizzes[title] = qid
	}
	jquiz, err := json.Marshal(quizzes)
	if writeErr(err, w) {
		return
	}
	fmt.Fprintf(w, string(jquiz)) //TODO WriteSuccess?
}
