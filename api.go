package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
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
		WARNING.Println("Create Class - User not authenticated")
		return
	}

	j := struct {
		Name     string              `json:"name"`
		Students []map[string]string `json:"students"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&j)
	if writeErr(err, w) {
		return
	}

	//create sid for each student
	for _, s := range j.Students {
		s = createStudentId(s)
	}

	st, err := json.Marshal(j.Students)
	if writeErr(err, w) {
		return
	}

	// insert the quiz
	var cid int
	err = db.QueryRow(`INSERT INTO classes (name, uid, students) 
		VALUES($1, $2, $3)  RETURNING cid`, j.Name, user.Uid, string(st)).Scan(&cid)
	if writeErr(err, w) {
		ERROR.Println("Create Class - INSERT uid=" + strconv.Itoa(user.Uid))
		return
	}
	TRACE.Println("Create Class - INSERT cid=" + strconv.Itoa(cid))

	//send student emails
	for _, s := range j.Students {
		go sendStudentClassEmail(cid, j.Name, s)
	}
	writeSuccess(w, cid)
}

func createStudentId(student map[string]string) map[string]string {
	b := make([]byte, 12)
	rand.Read(b)

	str := base64.StdEncoding.EncodeToString(b)
	//fmt.Println(str)
	student["sid"] = str
	return student
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
func handleAddStudent(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if user == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Add Student - User not authenticated")
		return
	}

	cid := mux.Vars(r)["cid"]

	decoder := json.NewDecoder(r.Body)
	j := struct {
		Cid   int    `json:"cid"`
		Email string `json:"email"`
		Fname string `json:"fname"`
		Lname string `json:"lname"`
	}{}
	err := decoder.Decode(&j)
	if writeErr(err, w) {
		return
	}

	student := make(map[string]string)
	student["email"] = j.Email
	student["fname"] = j.Fname
	student["lname"] = j.Lname

	//TODO ->> json
	var jstring, cname string
	err = db.QueryRow(`SELECT name, students FROM classes WHERE uid = $1 AND cid = $2`, user.Uid, cid).Scan(&cname, &jstring)
	if writeErr(err, w) {
		ERROR.Println("Add Student - SELECT cid=" + cid)
		return
	}

	// create student id
	student = createStudentId(student)

	var students []map[string]string
	err = json.Unmarshal([]byte(jstring), &students)
	if writeErr(err, w) {
		return
	}
	//students[j.Sid] = j.Name
	students = append(students, student)

	js, err := json.Marshal(students)
	if writeErr(err, w) {
		return
	}

	_, err = db.Exec(`UPDATE classes SET students=$1 WHERE cid= $2`, string(js), cid)
	if writeErr(err, w) {
		ERROR.Println("Add Student - UPDATE cid=" + cid)
		return
	}
	TRACE.Println("Add Student - UPDATE cid=" + cid)

	go sendStudentClassEmail(j.Cid, cname, student)

	fmt.Println("TESTING: ", student)

	writeSuccess(w, student)
}

// URL: /classes/{cid}/students
//
//Expecting JSON :
// {
//   "sid": string
// }
func handleDeleteStudent(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Delete Student - User not authenticated")
		return
	}
	cid := mux.Vars(r)["cid"]

	decoder := json.NewDecoder(r.Body)
	id := struct {
		Sid string `json:"sid"`
	}{}
	err := decoder.Decode(&id)
	if writeErr(err, w) {
		return
	}

	var jstring, cname string
	err = db.QueryRow(`SELECT name, students 
                     FROM classes 
                     WHERE uid = $1 
                     AND cid = $2`, user.Uid, cid).Scan(&cname, &jstring)
	if writeErr(err, w) {
		ERROR.Println("Delete Student - SELECT cid=" + cid)
		return
	}

	var students []map[string]string
	err = json.Unmarshal([]byte(jstring), &students)
	if writeErr(err, w) {
		return
	}

	//fmt.Println("id.SID ", id.Sid)
	var newStudents []map[string]string
	for i, _ := range students {
		//fmt.Println("\nsid ", students[i]["sid"])
		if !strings.EqualFold(students[i]["sid"], id.Sid) {
			//fmt.Println(" !=  ", students[i]["sid"])
			//append(newStudents, students[student])
			newStudents = append(newStudents, students[i])
		}
	}

	//fmt.Println("\n\nNEW ", newStudents)
	js, err := json.Marshal(newStudents)
	if writeErr(err, w) {
		return
	}

	_, err = db.Exec(`UPDATE classes 
                    SET students=$1
                    WHERE cid=$2`, string(js), cid)
	if writeErr(err, w) {
		ERROR.Println("Delete Student - UPDATE cid=" + cid)
		return
	}
	TRACE.Println("Add Student - UPDATE cid=" + cid)

	writeSuccess(w)

}

// URL: /classes/{cid}/students
//
//Expecting JSON :
// {
//   "sid": string
//   "email": string
// }
//TODO make this a little nicer..
func handleUpdateStudent(w http.ResponseWriter, r *http.Request) {
	user := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Update Student - User not authenticated")
		return
	}
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if writeErr(err, w) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	id := struct {
		Sid   string `json:"sid"`
		Email string `json:"email"`
	}{}
	err = decoder.Decode(&id)
	if writeErr(err, w) {
		return
	}

	var jstring, cname string
	err = db.QueryRow(`SELECT name, students 
                     FROM classes 
                     WHERE uid = $1 
                     AND cid = $2`, user.Uid, cid).Scan(&cname, &jstring)
	if writeErr(err, w) {
		ERROR.Println("Update Student - SELECT cid=", cid)
		return
	}

	var students []map[string]string
	err = json.Unmarshal([]byte(jstring), &students)
	if writeErr(err, w) {
		return
	}

	//fmt.Println("id.SID ", id.Sid)
	var newStudents []map[string]string
	var fname string
	var lname string
	for i, _ := range students {
		//fmt.Println("\nsid ", students[i]["sid"])
		if strings.EqualFold(students[i]["sid"], id.Sid) {
			//fmt.Println("\nEQUAL", students[i]["sid"], " ", id.Sid)
			students[i]["email"] = id.Email
			fname = students[i]["fname"]
			lname = students[i]["lname"]
			newStudents = append(newStudents, students[i])
		} else {
			//fmt.Println("\nNOT EQUAL ", students[i]["sid"])
			//append(newStudents, students[student])
			newStudents = append(newStudents, students[i])
		}
	}

	//fmt.Println("\n\nNEW ", newStudents)
	js, err := json.Marshal(newStudents)
	if writeErr(err, w) {
		return
	}

	_, err = db.Exec(`UPDATE classes 
                    SET students=$1
                    WHERE cid=$2`, string(js), cid)
	if writeErr(err, w) {
		ERROR.Println("Update Student - UPDATE cid=", cid)
		return
	}
	TRACE.Println("Update Student - UPDATE cid=", cid)

	student := make(map[string]string)
	student["email"] = id.Email
	student["fname"] = fname
	student["lname"] = lname
	student["sid"] = id.Sid

	go sendStudentClassEmail(cid, cname, student)

	writeSuccess(w)

}

func handleClassGet(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		WARNING.Println("Get Class - User not authenticated")
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
		ERROR.Println("Get Class - SELECT cid=" + cid)
		return
	}

	c := struct {
		Name     string              `json:"name"`
		Students []map[string]string `json:"students"`
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
		WARNING.Println("Get Class List - User not authenticated")
		return
	}

	rows, err := db.Query(`
    SELECT cid, name
    FROM classes
    WHERE classes.uid = $1
  `, auth.Uid)
	if writeErr(err, w) {
		ERROR.Println("Get Class List - SELECT uid=" + strconv.Itoa(auth.Uid))
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

//can only change class name right now.. anything else we want to change?
//TODO: check this
//Expecting JSON body of the form
//{
//  "name":string
//}
func handleClassUpdate(w http.ResponseWriter, r *http.Request) {

	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		return
	}

	vars := mux.Vars(r)
	cid, err1 := strconv.Atoi(vars["cid"])
	if writeErr(err1, w) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	j := struct {
		Name string `json:"name"`
	}{}
	err := decoder.Decode(&j)
	if writeErr(err, w) {
		return
	} else {
		_, err := db.Exec(`UPDATE classes SET name=$1 WHERE cid=$2`, j.Name, cid)
		if writeErr(err, w) {
			fmt.Println("\nDB err")
			return
		} else {
			writeSuccess(w)
		}
	}
}

func handleClassDelete(w http.ResponseWriter, r *http.Request) {
	auth := auth(r)
	if auth == nil {
		writeErr(fmt.Errorf("User not authenticated"), w)
		return
	}
	vars := mux.Vars(r)
	cid, err := strconv.Atoi(vars["cid"])
	if writeErr(err, w) {
		return
	}

	_, err = db.Exec(`DELETE FROM quiz 
                    WHERE cid=$1`, cid)
	if writeErr(err, w) {
		return
	}

	_, err = db.Exec(`DELETE FROM classes 
                    WHERE cid=$1`, cid)
	if writeErr(err, w) {
		return
	} else {
		writeSuccess(w)
	}
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
    INSERT INTO quiz (info, cid)
		VALUES($1, $2) RETURNING qid
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

//Expecting
//    {
//      grades : {
//        studentid (int) : grade (int),
//        ...
//      }
//    }
//
//experimenting
//func handleGradeList(w http.ResponseWriter, r *http.Request) {
//fmt.Println("here")

//auth := auth(r)
//if auth == nil {
//writeErr(fmt.Errorf("User not authenticated"), w)
//WARNING.Println("Get Grade List - User not authenticated")
//return
//}

//cid := mux.Vars(r)["id"]

//g := struct {
//Sid   int `json:"studentid"`
//Grade int `json:"grade"`
//}{}

//err := json.NewDecoder(r.Body).Decode(&g)
//if writeErr(err, w) {
//return
//}

////find name for each student
//rows, err := db.Query(`
//SELECT students->>'fname'
//FROM classes
//WHERE classes.cid = $1
//AND students->>'sid' = $2`, cid, g.Sid)
//if writeErr(err, w) {
//ERROR.Println("Get Grade List - SELECT")
//return
//}

//defer rows.Close()
//for rows.Next() {
//fmt.Println(rows)
//}

//}

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
      SELECT qid, info->>'title', name
      FROM quiz, classes
      WHERE classes.uid = $1
      AND classes.cid = $2
      AND quiz.cid = classes.cid
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
		ERROR.Println("Get Quiz List - SELECT")
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
