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
  //"strconv"
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
	//TODO drop tables
}

func homePage(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//TODO ...see a pattern here for *Page(w, r)
func registerPage(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email, pass := r.FormValue("email"), r.FormValue("password")
	fmt.Fprintf(w, "<html><body><h1>Hello, %s %s!</h1></body></html>", email, pass)
	//TODO check for authenticiousnessity
	
	
	var userid int
	var pw string
	err := db.QueryRow("SELECT user_id, password FROM \"Users\" WHERE username=$1", email).Scan(&userid, &pw)
    switch {
    case err == sql.ErrNoRows:
            fmt.Fprintf(w, "No user with that ID.")
    case err != nil:
            http.Error(w, err.Error(), http.StatusInternalServerError)
    default:
            fmt.Fprintf(w, "user_id is %b \n\n password is %s", userid, pw)
    }
    
    
	//TODO return site wide cookie & add to DB
}

func register(w http.ResponseWriter, r *http.Request) {
	email, pass := r.FormValue("email"), r.FormValue("password")

	salt := make([]byte, 32)
	_, err := rand.Read(salt)

	//FIXME just to see if this worked
	fmt.Fprintf(w, "%s %s %s", email, pass, salt)

	//create secure, salted hash
	p := append([]byte(pass), salt...)
	hash := sha256.Sum256(p)

	//TODO check existence?
	_, err = db.Exec(`INSERT INTO teachers(email, password, salt)
    VALUES($1, $2, $3)`, email, hash, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return //necessary?
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	//FIXME stub to check for cookie
}

func handleAnswer(s http.ResponseWriter, r *http.Request){
 // vars := mux.Vars(r)
 // qID, err := strconv.Atoi(vars["id"])
 // if err != nil {
 //   fmt.Println(err)
 // } else {
 //   fmt.Printf("%d", qID) //testing
 // }
 fmt.Println("here")
}

func main() {
	r := mux.NewRouter()
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/static/"))))
	r.HandleFunc("/", homePage)
	r.HandleFunc("/login", loginPage).Methods("GET")
	r.HandleFunc("/logmein", login).Methods("POST")
	r.HandleFunc("/register", registerPage).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")

	//TODO these are just ideas
	//r.HandleFunc("/quiz/{id}", handleQuizGet).Methods("GET")
	r.HandleFunc("/quiz/{id}", handleAnswer).Methods("PUT")
	//r.HandleFunc("/quiz/{id}/edit", handleQuizEdit).Methods("POST, GET")
	//r.HandleFunc("/quiz/add, handleQuizCreate).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
