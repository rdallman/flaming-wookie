package main

import (
	"crypto/rand"
	"encoding/base64"
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

	str := base64.URLEncoding.EncodeToString(b)
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

	writeSuccess(w, student)
}

// URL: /classes/{cid}/students
//
// Expecting
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

	student := struct {
		Sid string `json:"sid"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&student)
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

	for i, stud := range students {
		if stud["sid"] == student.Sid {
			students = append(students[:i], students[i+1:]...)
			break
		}
	}

	js, err := json.Marshal(students)
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
// Expecting:
// {
//   "sid": string
//   "fname": string
//   "lname": string
//   "email": string
// }
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

	student := struct {
		Sid   string `json:"sid"`
		Fname string `json:"fname"`
		Lname string `json:"lname"`
		Email string `json:"email"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&student)
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

	for i, stud := range students {
		if stud["sid"] == student.Sid {
			if student.Email != "" {
				students[i]["email"] = student.Email
			}
			if student.Fname != "" {
				students[i]["fname"] = student.Fname
			}
			if student.Lname != "" {
				students[i]["lname"] = student.Lname
			}
			break
		}
	}

	js, err := json.Marshal(students)
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

	var classes []map[string]interface{}
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
		WARNING.Println("Delete Class - User not authenticated")
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
		ERROR.Println("Delete Class - DELETE FROM quiz WHERE cid=" + strconv.Itoa(cid))
		return
	}

	_, err = db.Exec(`DELETE FROM classes 
                    WHERE cid=$1`, cid)
	if writeErr(err, w) {
		ERROR.Println("Delete Class - DELETE cid=" + strconv.Itoa(cid))
		return
	} else {
		writeSuccess(w)
	}
	TRACE.Println("Delete Class - DELETE cid=" + strconv.Itoa(cid))
}
