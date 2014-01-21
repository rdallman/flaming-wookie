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
	"time"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
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

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://wookie:password@absker.com/wookie?sslmode=disable")
	if err != nil {
		panic(err)
	}
	//TODO NEED TO USE BYTEA TO STORE HASHES
	//_, err = db.Exec(`ALTER TABLE "Users" ALTER COLUMN password TYPE bytea`)
	//fmt.Println(err)
	//_, err = db.Exec(`ALTER TABLE "Users" ALTER COLUMN salt TYPE bytea`)
	//fmt.Println(err)
	// TODO drop tables
	// DANGER this will empty the db
	// I'm going to do this just to change the names because it's a pain in the ass
	// and nice for testing (and reading types and such, to boot)

	//_, err = db.Exec(`DROP TABLE "Classes" "Users" "Students"`)
	//fmt.Println(err)

	//rows, err := db.Query("SELECT tablename from pg_catalog.pg_tables")
	//for rows.Next() {
	//var tablename string
	//rows.Scan(&tablename)
	//fmt.Println(tablename)
	//}
	//_, err = db.Exec(`CREATE TABLE teachers
	//(tid integer,
	//email text,
	//password text,
	//salt text)`)

	//_, err = db.Exec(`CREATE TABLE classes
	//(cid integer,
	//name text,
	//active boolean,
	//uid text)`) //TODO not sure about this one?

	////TODO this just feels wrong
	//_, err = db.Exec(`CREATE TABLE students
	//(sid text,
	//pin integer,
	//cid integer)`)

	//_, err = db.Exec(`CREATE TABLE quiz
	//(qid integer,
	//name text,
	//info json)`)

	//TODO solution for grading
}

func handlePage(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth(w, r)
		if err == nil {
			// TODO redirect to dashboard
		}

		err = templates.ExecuteTemplate(w, name+".html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	inputEmail, inputPass := r.FormValue("email"), r.FormValue("password")
	fmt.Fprintf(w, inputPass, inputEmail)

	//dbpw is the salted sha256 hash we stored as password
	var salt, dbpw string
	err := db.QueryRow(`SELECT salt, password FROM "Users" WHERE username=$1`, inputEmail).Scan(&salt, &dbpw)

	//salt input password, hash and compare to database salted hash
	hash := sha256.Sum256(append([]byte(inputPass), salt...))
	phash := string(hash[:]) //finicky

	switch {
	case err == sql.ErrNoRows:
		// TODO add flash messages
		fmt.Fprintf(w, "No user with that username.")
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case dbpw != phash:
		// TODO add flash messages
		fmt.Fprintf(w, "Incorrect password")
	default:
		expire := time.Now().AddDate(0, 1, 0)
		cookie := &http.Cookie{
			Name:    "logged-in",
			Value:   inputEmail,
			Expires: expire,
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 307)
	}
	// TODO return site wide cookie & add to DB
}

func register(w http.ResponseWriter, r *http.Request) {
	email, pass := r.FormValue("email"), r.FormValue("password")

	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	//saltstr := string(salt[:])

	// create secure, salted hash
	hash := sha256.Sum256(append([]byte(pass), salt...))
	phash := string(hash[:])

	// TODO check existence?
	_, err = db.Exec(`INSERT INTO "Users" (username, password, salt)
    VALUES($1, $2, $3)`, email, phash, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/login", 307)
}

// checks for cookie; if no cookie then redirect to home page
//
// this function... I don't think it does what you think it does
// I am Inigo Montoya. You kill my father. Prepare to die.
//
// should be done on every request (sigh, I know) and display different options for
// teachers, alternatively you could return a user each time and if not nil, you have a teacher
// then with that knowledge, get teacher specific data from database (in other methods...).
// going off of that, I guess you really just need to return an ID each time since that'll be
// what's used to hit the DB. That may prove problematic but premature optimization is the root of all... well, you know
func auth(w http.ResponseWriter, r *http.Request) error {
	// FIXME stub to check for cookie
	_, err := r.Cookie("logged-in")
	if err != nil {
		return err
	} else {
		// fmt.Fprintf(w, "You has cookie: %s", cookie.Value)
		return nil
	}
}

func main() {
	r := mux.NewRouter()
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/static/"))))
	r.HandleFunc("/", handlePage("index"))
	r.HandleFunc("/login", handlePage("login")).Methods("GET")
	r.HandleFunc("/logmein", login).Methods("POST")
	r.HandleFunc("/register", handlePage("register")).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")

	// TODO these are just ideas
	// r.HandleFunc("/quiz/{id}", handleQuizGet).Methods("GET")
	// r.HandleFunc("/quiz/{id}", handleAnswer).Methods("PUT")
	// r.HandleFunc("/quiz/{id}/edit, handleQuizEdit).Methods("POST, GET")
	// r.HandleFunc("/quiz/add, handleQuizCreate).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
