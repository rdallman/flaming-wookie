package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"

  	"strconv"

	"time"

)

var templates = template.Must(template.New("").Delims("<<<", ">>>").ParseGlob("templates/*.html"))
var db *sql.DB

type Quiz struct {
	Questions []Question
}

type Question struct {
	Text    string
	Answers []Answer
}

type Answer struct {
	Text    string
	Correct bool
}

type User struct{
	Uid int
	Email string
}

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://wookie:password@absker.com/wookie?sslmode=disable")
	//db, err = sql.Open("postgres", "user=reed dbname=wookie sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	//////////////////////////////////////
	// drop tables
	// DANGER this will empty the db
	//
	//////////////////////////////////////
	//
	//_, err = db.Exec(`DROP TABLE "Classes", "Users", "Students" CASCADE`)
	//fmt.Println(err)
	//_, err = db.Exec(`DROP TABLE classes, users, students, quiz, class_student CASCADE`)
	//fmt.Println(err)

	//for getting table names -- handy
	//
	//rows, err := db.Query("SELECT tablename from pg_catalog.pg_tables")
	//for rows.Next() {
	//var tablename string
	//rows.Scan(&tablename)
	//fmt.Println(tablename)
	//}

	/////////////////////////////////////////////
	//creating
	//TODO NOT NULL all of these later...
	/////////////////////////////////////////////

	//_, err = db.Exec(`CREATE TABLE users (
	//uid serial PRIMARY KEY,
	//email text UNIQUE,
	//password bytea,
	//salt bytea
	//)`)
	//fmt.Println(err)

	//_, err = db.Exec(`CREATE TABLE classes (
	//cid serial PRIMARY KEY,
	//name text,
	//uid integer REFERENCES users (uid)
	//)`)
	//fmt.Println(err)

	////TODO this just feels wrong
	//_, err = db.Exec(`CREATE TABLE students (
	//sid serial PRIMARY KEY,
	//schoolid text,
	//pin integer
	//)`)
	//fmt.Println(err)

	//_, err = db.Exec(`CREATE TABLE quiz (
	//qid serial PRIMARY KEY,
	//title text,
	//info json,
	//cid integer REFERENCES classes (cid)
	//)`)
	//fmt.Println(err)

	////for authentication of students
	//_, err = db.Exec(`CREATE TABLE class_student (
	//cid integer REFERENCES classes (cid),
	//sid integer REFERENCES students (sid)
	//)`)
	//fmt.Println(err)

	////////////////////////////////
	//TODO FIXME STAHP OTHER KEYWORDS
	//WIP
	///////////////////////////////

	//TODO some thought needed... this would be a shitton of rows
	//_, err = db.Exec(`CREATE TABLE quiz_student_question_answer (
	//qid integer REFERENCES quiz (qid),
	//sid integer REFERENCES student (sid),
	//number integer,
	//answer text
	//)`)

	//TODO solution for grading
}

func handlePage(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, email := auth(w, r) //check for cookie
		var err error
		if uid == -1 { //no cookie
			err = templates.ExecuteTemplate(w, name+".html", nil)
		} else { //cookie - execute template with 'User' struct
			user := &User{Uid: uid, Email: email}
			err = templates.ExecuteTemplate(w, name+".html", user)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	inputEmail, inputPass := r.FormValue("email"), r.FormValue("password")

	//dbpw is the salted sha256 hash we stored as password
	var salt, dbpw string
	err := db.QueryRow(`SELECT salt, password FROM users WHERE email=$1`, inputEmail).Scan(&salt, &dbpw)

	//salt input password, hash and compare to database salted hash
	hash := sha256.Sum256(append([]byte(inputPass), salt...))
	phash := string(hash[:]) //finicky

	switch {
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case err == sql.ErrNoRows || dbpw != phash:
		// TODO add flash messages
		fmt.Fprintf(w, "Invalid login credentials")
	default:
		createCookie(w, inputEmail)
		//TODO this is weird (flash)
		http.Redirect(w, r, "/", 307)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	email, pass1, pass2 := r.FormValue("email"), r.FormValue("password"), r.FormValue("password confirm")

	//check if any fields are blank or if passwords do not match
	if email=="" || pass1=="" || pass2=="" || pass1!=pass2 {
		fmt.Fprintf(w, "Invalid information")
		return //do not add to db
	} 
	//check to see if email already exists in db
	var id int
	er := db.QueryRow(`SELECT uid FROM users WHERE email=$1`, email).Scan(&id)
	if er != sql.ErrNoRows {
		fmt.Fprintf(w, "User already exists")
		return //do not add to db
	}

	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	//saltstr := string(salt[:])

	// create secure, salted hash
	hash := sha256.Sum256(append([]byte(pass1), salt...))
	phash := string(hash[:])

	_, err = db.Exec(`INSERT INTO users (email, password, salt)
    VALUES($1, $2, $3)`, email, phash, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	//(login) create cookie and redirect to homepage
	createCookie(w, email)
	http.Redirect(w, r, "/", 307)
}

//deletes cookie and redirects to homepage
func logout(w http.ResponseWriter, r *http.Request) {
	expire := time.Now()
	cookie := &http.Cookie{
		Name:    "logged-in",
		Value:   "",
		Expires: expire,
		Path:    "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 307)
	//TODO (when we add cookies to db) delete cookie from db
}

//creates and sets cookie for user
func createCookie(w http.ResponseWriter, email string) {
	expire := time.Now().AddDate(0, 1, 0)
	cookie := &http.Cookie{
		Name:    "logged-in",
		Value:   email,
		Expires: expire,
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	//TODO add cookie to db
}

//checks for cookie
// 	if cookie -> returns the user's uid and email
//	if no cookie (or invalid) -> return -1 and ""
func auth(w http.ResponseWriter, r *http.Request) (int, string) {
	cookie, err := r.Cookie("logged-in")
	if err != nil { // no cookie
		return -1, ""
	} else { //cookie
		var uid int
		err = db.QueryRow(`SELECT uid FROM users WHERE email=$1`, cookie.Value).Scan(&uid)
		if err != nil { // invalid user
			return -1, ""
		}
		return uid, cookie.Value // valid user
	}
}

// add more later
func handleQuizGet(w http.ResponseWriter, r *http.Request) {  
  vars := mux.Vars(r)
  qID, err := strconv.Atoi(vars["id"])
  if err != nil {
    fmt.Println(err)
  } else {
    //fmt.Printf("%d", qID) //testing 
    rows, err := db.Query(`SELECT * FROM quiz WHERE qid=$1`, qID)
    if err != nil {
    fmt.Printf("%s", err)
    } else { 
      for rows.Next() {
        //fmt.Printf("here") //testing
        var qid int
        var title string
        var info string
        var cid int
        err = rows.Scan(&qid, &title, &info, &cid) 
        if err != nil {
        fmt.Printf("%s", err)
        } else {
          fmt.Printf("\nqid:%d \ttitle:%s \tinfo:%s \tcid:%d", qid, title, info, cid)
      }
      }
    }
}
} 

func main() {
	r := mux.NewRouter()
	// r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.HandleFunc("/", handlePage("index"))
	r.HandleFunc("/login", handlePage("login")).Methods("GET")
	r.HandleFunc("/logmein", login).Methods("POST")
	r.HandleFunc("/logmeout", logout)
	r.HandleFunc("/register", handlePage("register")).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")


	//TODO these are just ideas
	r.HandleFunc("/quiz/{id}", handleQuizGet).Methods("GET")
	//r.HandleFunc("/quiz/{id}", handleAnswer).Methods("PUT")
	//r.HandleFunc("/quiz/{id}/edit", handleQuizEdit).Methods("POST, GET")
	//r.HandleFunc("/quiz/add, handleQuizCreate).Methods("POST")

	r.HandleFunc("/dashboard", handlePage("dashboard")).Methods("GET")
	r.HandleFunc("/quiz", handlePage("quiz")).Methods("GET")

	// TODO these are just ideas
	// r.HandleFunc("/quiz/{id}", handleQuizGet).Methods("GET")
	// r.HandleFunc("/quiz/{id}", handleAnswer).Methods("PUT")
	// r.HandleFunc("/quiz/{id}/edit, handleQuizEdit).Methods("POST, GET")
	r.HandleFunc("/quiz/add", handleQuizCreate).Methods("POST")


	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
