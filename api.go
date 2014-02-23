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
		writeErr(fmt.Errorf("User not authenticated"), w)
		return
	}

	//TODO generate ids here...

	j := struct {
		Name     string              `json:"name"`
		Students []map[string]string `json:"students"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&j)
	if writeErr(err, w) {
		return
	}
	fmt.Println(j)
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

// TODO needs tidying, put in URL? JSON?
// TODO just use /class UPDATE method?
//
// add a student to a class
//
// URL: /classes/{cid}/students
//
// expecting JSON body:
//  {
//    "name": string,
//    "sid": string
//  }
func handleAddStudents(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		return
	}

	cid := mux.Vars(r)["cid"]

	decoder := json.NewDecoder(r.Body)
	j := struct {
		Name string `json:"name"`
		Sid  string `json:"sid"`
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
	students[j.Sid] = j.Name

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

func handleClassGet(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		return
	}

	cid := mux.Vars(r)["cid"]
	var name, students string

	err := db.QueryRow(`
    SELECT  name, students
    FROM classes
    WHERE classes.cid = $1
  `, cid).Scan(&name, &students)

	if writeErr(err, w) {
		return
	}

	c := struct {
		Name     string                   `json:"name"`
		Students []map[string]interface{} `json:"students"`
	}{
		name,
		nil,
	}

	err = json.Unmarshal([]byte(students), &c.Students)
	if writeErr(err, w) {
		return
	}
	writeSuccess(w, c)
}

// GET /classes
func handleClassList(w http.ResponseWriter, r *http.Request) {
	//title and id, return JSON
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		return
	}

	rows, err := db.Query(`
    SELECT cid, name
    FROM classes
    WHERE classes.uid = $1
  `, auth.Uid)
	if writeErr(err, w) {
		return
	}
	defer rows.Close()

	classes := make([]map[string]interface{}, 0)
	for rows.Next() {
		var cid int
		var name string
		err = rows.Scan(&cid, &name)
		if writeErr(err, w) {
			return
		}
		classes = append(classes, map[string]interface{}{"name": name, "cid": cid})
	}
	writeSuccess(w, classes)
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

// handleQuizCreate creates a quiz from an AJAX POST request
func handleQuizCreate(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
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
    `, string(info), cid)
	if writeErr(err, w) {
		return
	}
	writeSuccess(w)
}

// handleQuizGet qets the quizID from the end of the given URL w
// and writes the info json from the db back using r.
// ((add more later))

//TODO auth?
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
	var obj map[string]interface{}
	json.Unmarshal([]byte(info), &obj)
	writeSuccess(w, obj)
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

// TODO Delete a quiz
func handleQuizDelete(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
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

// GET /quiz
// GET /classes/{cid}/quiz
func handleQuizList(w http.ResponseWriter, r *http.Request) {
	//title and id, return JSON
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		return
	}

	vars := mux.Vars(r)

	var rows *sql.Rows
	var err error

	//TODO maybe fall through with string.. handling of rows is bad
	if cid, ok := vars["cid"]; ok {
		rows, err = db.Query(`
      SELECT qid, info->>'title', name
      FROM quiz, classes
      WHERE classes.uid = $1
      AND classes.cid = $2
      `, auth.Uid, cid)
	} else {
		rows, err = db.Query(`
      SELECT qid, info->>'title', name
      FROM quiz, classes
      WHERE classes.uid = $1
      AND classes.cid = quiz.cid
    `, auth.Uid)
	}
	if writeErr(err, w) {
		return
	}
	defer rows.Close() //TODO these may not close if err != sql.NoRowsErr

	qs := make([]map[string]interface{}, 0)

	for rows.Next() {
		var qid int
		var title, name string
		err = rows.Scan(&qid, &title, &name)
		if writeErr(err, w) {
			return
		}
		qs = append(qs, map[string]interface{}{"title": title, "qid": qid, "name": name})
	}
	writeSuccess(w, qs)
}
