package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
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
	decoder := json.NewDecoder(r.Body)

	j := struct {
		name     string
		students map[string]string //map[id]name
	}{}

	err := decoder.Decode(&j)
	if writeErr(err, w) {
		return
	}

	st, err := json.Marshal(j.students)

	// insert the quiz
	_, err = db.Exec(`INSERT INTO classes (name, uid, students) 
		VALUES($1, $2, $3)`, j.name, user.Uid, string(st))
	if writeErr(err, w) {
		return
	}
	writeSuccess(w)
}

func handleAddStudents(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	j := struct {
		name string
		sid  string
		cid  int
	}{}
	err := decoder.Decode(&j)
	if writeErr(err, w) {
		return
	}

	var jstring string
	err = db.QueryRow(`SELECT students FROM classes WHERE uid = $1 AND cid = $2`, user.Uid, j.cid).Scan(&jstring)
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

	_, err = db.Exec(`UPDATE classes SET students=$1 WHERE cid= $2`, string(js), j.cid)
	if writeErr(err, w) {
		return
	}

	writeSuccess(w)
}

// TODO: flash message to show quiz was added, and redirect
//
// Expecting http body with JSON of form:
//  cid : int
//  title : string
//  info :
//    questions : [
//      {
//        text : string,
//        answers : [
//          string,
//          ...
//        ],
//        correct : string
//      },
//      ...
//    ],
//    grades : {
//      studentid (int) : grade (int)
//    }
//
// on creation just make a blank map for grades

// TODO add cid to frontend

// handleQuizCreate creates a quiz from an AJAX POST request
func handleQuizCreate(w http.ResponseWriter, r *http.Request) {
	// grab body of request (should be the json of the quiz)
	decoder := json.NewDecoder(r.Body)

	j := struct {
		cid   int
		title string
		info  Quiz
	}{}

	err := decoder.Decode(&j)
	if writeErr(err, w) {
		return
	}

	// insert the quiz
	_, err = db.Exec(`INSERT INTO quiz (title, info, cid) 
		VALUES($1, $2, $3)`, j.title, j.info, j.cid)
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
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Printf("%d", qID) //testing
		var info string
		err := db.QueryRow(`SELECT info FROM quiz WHERE qid=$1`, qID).Scan(&info)
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			fmt.Fprintf(w, info)
		}
	}
}

// handleQuizUpdate qets the quizID from the end of the given URL w,
// gets the form data, and updates the quiz in the db.
// ((just an idea, not sure if we actually need this))
func handleQuizUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HERE")
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

// handleQuizList returns a json of all quiz ids and titles
// using r.
func handleQuizDelete(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth != nil {
		vars := mux.Vars(r)
		qid, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(qid)
			_, err := db.Exec(`DELETE FROM quiz WHERE qid=$1`, qid)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("\nDeleted.")
			}
		}
	} else { //this is bad, but we can decide this later...
		fmt.Println("Error - you cannot delete a quiz")
	}
}

func handleQuizList(w http.ResponseWriter, r *http.Request) {
	//title and id, return JSON
	rows, err := db.Query(`SELECT qid, title FROM quiz`)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		quizzes := make(map[string]int)
		for rows.Next() {
			var qid int
			var title string
			err = rows.Scan(&qid, &title)
			if err != nil {
				fmt.Println(err)
			} else {
				quizzes[title] = qid
			}
		}
		jquiz, err := json.Marshal(quizzes)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Fprintf(w, string(jquiz))
		}
	}
}
