package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

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
		WARNING.Println("Create Quiz - User not authenticated")
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
	//_, err = db.Exec(`
	//INSERT INTO quiz (info, cid)
	//VALUES($1, $2)
	//`, string(info), cid)
	var qid int
	err = db.QueryRow(`
    INSERT INTO quiz (info, cid, type)
		VALUES($1, $2, 1) RETURNING qid
    `, string(info), cid).Scan(&qid)
	if writeErr(err, w) {
		ERROR.Println("Create Quiz - INSERT cid=" + cid)
		return
	}
	TRACE.Println("Create Quiz - INSERT qid=" + strconv.Itoa(qid))
	writeSuccess(w, qid)
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
		ERROR.Println("Get Quiz - SELECT qid=" + strconv.Itoa(qID))
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
			ERROR.Println("Update Quiz - UPDATE qid=" + strconv.Itoa(qid))
		} else {
			TRACE.Println("Update Quiz - UPDATE qid=" + strconv.Itoa(qid))
		}
	}
}

// TODO Delete a quiz
func handleQuizDelete(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Delete Quiz - User not authenticated")
		return
	}
	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["id"])
	if writeErr(err, w) {
		return
	}
	_, err = db.Exec(`DELETE FROM quiz WHERE qid=$1`, qid)
	if writeErr(err, w) {
		ERROR.Println("Delete Quiz - DELETE qid=" + strconv.Itoa(qid))
		return
	}
	TRACE.Println("Delete Quiz - DELETE qid=" + strconv.Itoa(qid))
	writeSuccess(w)
}

// GET /quiz
// GET /classes/{cid}/quiz
func handleQuizList(w http.ResponseWriter, r *http.Request) {
	//title and id, return JSON
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Get Quiz List - User not authenticated")
		return
	}

	vars := mux.Vars(r)

	var rows *sql.Rows
	var err error

	//TODO maybe fall through with string.. handling of rows is bad
	if cid, ok := vars["cid"]; ok {
		rows, err = db.Query(`
      SELECT qid, info->>'title', info->>'grades', name
      FROM quiz, classes
      WHERE classes.uid = $1
      AND classes.cid = $2
      AND quiz.cid = classes.cid
      AND quiz.type = 1
      `, auth.Uid, cid)
	} else {
		rows, err = db.Query(`
      SELECT qid, info->>'title', info->>'grades', name
      FROM quiz, classes
      WHERE classes.uid = $1
      AND classes.cid = quiz.cid
      AND quiz.type = 1
      `, auth.Uid)
	}
	if writeErr(err, w) {
		ERROR.Println("Get Quiz List - SELECT")
		return
	}
	defer rows.Close() //TODO these may not close if err != sql.NoRowsErr

	qs := make([]map[string]interface{}, 0)

	for rows.Next() {
		var qid int
		var title, name, gradesString string
		var grades map[string]int
		err = rows.Scan(&qid, &title, &gradesString, &name)
		if writeErr(err, w) {
			return
		}
		err = json.Unmarshal([]byte(gradesString), &grades)

		qs = append(qs, map[string]interface{}{"title": title, "qid": qid, "name": name, "showGrades": len(grades) > 0})
	}
	writeSuccess(w, qs)
}

func handleGetQuizGrades(w http.ResponseWriter, r *http.Request) {
	//title and id, return JSON
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Get Quiz Grades - User not authenticated")
		return
	}

	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["id"])
	if writeErr(err, w) {
		return
	}

	//get cid and quizinfo from quiz
	var cid int
	var quizinfo string
	err = db.QueryRow(`SELECT cid, info FROM quiz WHERE qid=$1`, qid).Scan(&cid, &quizinfo)
	if writeErr(err, w) {
		ERROR.Println("Get Quiz Grades - SELECT qid=" + strconv.Itoa(qid))
		return
	}

	//get students from class
	var studentJson string
	err = db.QueryRow(`SELECT students FROM classes WHERE uid = $1 AND cid = $2`, auth.Uid, cid).Scan(&studentJson)
	if writeErr(err, w) {
		ERROR.Println("Get Quiz Grades - SELECT cid=" + strconv.Itoa(cid))
		return
	}

	//unmarshal quiz
	var quiz Quiz
	err = json.Unmarshal([]byte(quizinfo), &quiz)
	if writeErr(err, w) {
		ERROR.Println("Get Quiz Grades - unmarshal quiz=" + strconv.Itoa(qid))
		return
	}

	//unmarshal students
	var students []map[string]string
	err = json.Unmarshal([]byte(studentJson), &students)
	if writeErr(err, w) {
		ERROR.Println("Get Quiz Grades - unmarshal students=" + strconv.Itoa(cid))
		return
	}

	// check if the quiz was taken
	if len(quiz.Grades) == 0 {
		writeSuccess(w, map[string]bool{"found": false})
		return
	}
	//map student name with their grade
	min := 100
	max := 0
	total := 0
	count := 0
	var studentGrades = make(map[string]int)
	for _, student := range students {
		if grade, found := quiz.Grades[student["sid"]]; found {
			studentGrades[student["fname"]+" "+student["lname"]] = grade
			if grade < min {
				min = grade
			}
			if grade > max {
				max = grade
			}
			total += grade
			count += 1
		} else {
			studentGrades[student["fname"]+" "+student["lname"]] = -1
		}
	}

	var gradeReturn = make(map[string]interface{})
	gradeReturn["max"] = max
	gradeReturn["min"] = min
	gradeReturn["avg"] = total / count
	gradeReturn["studentGrades"] = studentGrades

	writeSuccess(w, gradeReturn)
}

// POLLS

// handlePollCreate creates a quiz from an AJAX POST request
func handlePollCreate(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Create Poll - User not authenticated")
		return
	}

	cid := mux.Vars(r)["cid"]

	var p Quiz
	err := json.NewDecoder(r.Body).Decode(&p)
	if writeErr(err, w) {
		fmt.Println("Here")
		return
	}

	info, err := json.Marshal(p)
	if writeErr(err, w) {
		return
	}

	// insert the quiz
	//_, err = db.Exec(`
	//INSERT INTO quiz (info, cid)
	//VALUES($1, $2)
	//`, string(info), cid)
	var pid int
	err = db.QueryRow(`
    INSERT INTO quiz (info, cid, type)
		VALUES($1, $2, 2) RETURNING qid
    `, string(info), cid).Scan(&pid)
	if writeErr(err, w) {
		ERROR.Println("Create Poll - INSERT cid=" + cid)
		return
	}
	TRACE.Println("Create Poll - INSERT qid=" + strconv.Itoa(pid))
	writeSuccess(w, pid)
}

/*
Get the poll with the provided pid
*/
func handlePollGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.Atoi(vars["id"])
	if writeErr(err, w) {
		return
	}
	//fmt.Printf("%d", qID) //testing
	var info string
	err = db.QueryRow(`SELECT info FROM quiz WHERE qid=$1`, pid).Scan(&info)
	if writeErr(err, w) {
		ERROR.Println("Get Poll - SELECT qid=" + strconv.Itoa(pid))
		return
	}

	var obj map[string]interface{}
	json.Unmarshal([]byte(info), &obj)
	writeSuccess(w, obj)
}

// TODO Delete a poll
func handlePollDelete(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Delete Quiz - User not authenticated")
		return
	}
	vars := mux.Vars(r)
	pid, err := strconv.Atoi(vars["id"])
	if writeErr(err, w) {
		return
	}
	_, err = db.Exec(`DELETE FROM quiz WHERE qid=$1`, pid)
	if writeErr(err, w) {
		ERROR.Println("Delete Poll - DELETE qid=" + strconv.Itoa(pid))
		return
	}
	TRACE.Println("Delete Poll - DELETE qid=" + strconv.Itoa(pid))
	writeSuccess(w)
}

//
func handlePollsList(w http.ResponseWriter, r *http.Request) {
	//title and id, return JSON
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Get Poll List - User not authenticated")
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
      AND quiz.cid = classes.cid
      AND quiz.type = 2
      `, auth.Uid, cid)
	} else {
		rows, err = db.Query(`
      SELECT qid, info->>'title', name
      FROM quiz, classes
      WHERE classes.uid = $1
      AND classes.cid = quiz.cid
      AND quiz.type = 2
      `, auth.Uid)
	}
	if writeErr(err, w) {
		ERROR.Println("Get Poll List - SELECT")
		return
	}
	defer rows.Close() //TODO these may not close if err != sql.NoRowsErr

	polls := make([]map[string]interface{}, 0)

	for rows.Next() {
		var pid int
		var title, name string
		err = rows.Scan(&pid, &title, &name)
		if writeErr(err, w) {
			fmt.Println("here")
			return
		}

		polls = append(polls, map[string]interface{}{"title": title, "qid": pid, "name": name})
	}
	writeSuccess(w, polls)
}

func handleAttendanceList(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Get Attendance List - User not authenticated")
		return
	}

	var rows *sql.Rows
	var err error

	cid := mux.Vars(r)["cid"]

	rows, err = db.Query(`
		select students, date_created
		from attendance
		where cid = $1
		order by date_created
		`, cid)

	if writeErr(err, w) {
		ERROR.Println("Get Attendance List - SELECT")
		return
	}
	defer rows.Close()

	attendance := make([]map[string]interface{}, 0)

	for rows.Next() {
		var students string
		var date time.Time
		err = rows.Scan(&students, &date)
		if writeErr(err, w) {
			fmt.Println(err)
			return
		}
		var studentsJson []map[string]int
		_ = json.Unmarshal([]byte(students), &studentsJson)
		attendance = append(attendance, map[string]interface{}{"date": date, "students": studentsJson})
	}
	writeSuccess(w, attendance)

}

func handlePollResults(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Get Attendance List - User not authenticated")
		return
	}

	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["id"])

	var info string
	err = db.QueryRow(`SELECT by_question FROM session_dump WHERE qid=$1 order by date_created limit 1`, qid).Scan(&info)
	if writeErr(err, w) {
		ERROR.Println("Poll Results - get qid=" + strconv.Itoa(qid))
		return
	}
	var obj []map[string]interface{}
	json.Unmarshal([]byte(info), &obj)
	writeSuccess(w, obj)
}
