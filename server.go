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
	// TODO drop tables
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
	r.ParseForm()
	inputEmail, inputPass := r.FormValue("email"), r.FormValue("password")
	// TODO check for authenticiousnessity
	var pw string
	err := db.QueryRow("SELECT password FROM \"Users\" WHERE username=$1", inputEmail).Scan(&pw)
	switch {
	case err == sql.ErrNoRows:
		// TODO add flash messages
		fmt.Fprintf(w, "No user with that username.")
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case inputPass != pw:
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

	// FIXME just to see if this worked
	fmt.Fprintf(w, "%s %s %s", email, pass, salt)

	// create secure, salted hash
	p := append([]byte(pass), salt...)
	hash := sha256.Sum256(p)

	// TODO check existence?
	_, err = db.Exec(`INSERT INTO teachers(email, password, salt)
    VALUES($1, $2, $3)`, email, hash, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // necessary?
	}
}

// checks for cookie; if no cookie then redirect to home page
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
