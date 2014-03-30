package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	templates = template.Must(template.New("").Delims("<<<", ">>>").ParseGlob("templates/*.html"))
	db        *sql.DB
	qzSesh    = make(map[int]Session)
)

//student sends: ID:PIN as USER in header...

// handlePage renders the page template for the given name path.
// If a user is logged in (determined by calling auth(r)) the template
// is executed with the user's information.
// TODO figure out somewhere better to put this...
func handlePage(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth(r) //check for cookie
		err := templates.ExecuteTemplate(w, name+".html", user)
		if err != nil {
			ERROR.Println("Handle Page -", name, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func serveFile(url string, filename string) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
	r.HandleFunc("/", handlePage("index"))

	//Authentication (auth.go) //TODO javascript me
	r.HandleFunc("/login", handlePage("login")).Methods("GET")
	r.HandleFunc("/logmein", login).Methods("POST") //TODO /login
	r.HandleFunc("/logmeout", logout)
	r.HandleFunc("/register", handlePage("register")).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")

	//API class methods (api.go)
	r.HandleFunc("/classes", handleClassList).Methods("GET")
	r.HandleFunc("/classes", handleCreateClass).Methods("POST")
	r.HandleFunc("/classes/{cid:[0-9]+}", handleClassGet).Methods("GET")
	//TODO r.HandleFunc("/classes/{cid:[0-9]+}", handleClassUpdate).Methods("POST")
	//TODO r.HandleFunc("/classes/{cid:[0-9]+}", handleClassDelete).Methods("DELETE")
	r.HandleFunc("/classes/{cid:[0-9]+}/quiz", handleQuizList).Methods("GET")
	r.HandleFunc("/classes/{cid:[0-9]+}/quiz", handleQuizCreate).Methods("POST")
	r.HandleFunc("/classes/{cid:[0-9]+}/student", handleAddStudent).Methods("POST")
	//TODO /classes/{cid:[0-9]+}/students POST, DELETE (i.e. update student, delete student) //Do we want this? just update class?

	//API quiz methods (api.go)
	r.HandleFunc("/quiz", handleQuizList).Methods("GET")
	r.HandleFunc("/quiz/{id:[0-9]+}", handleQuizGet).Methods("GET")
	r.HandleFunc("/quiz/{id:[0-9]+}", handleQuizDelete).Methods("DELETE")
	//TODO r.HandleFunc("/quiz/{id:[0-9]+}", handleQuizUpdate).Methods("POST")

	//API quiz session (session.go)
	r.HandleFunc("/quiz/{id:[0-9]+}/state", changeState).Methods("PUT")
	r.HandleFunc("/quiz/{id:[0-9]+}/answer", handleAnswer).Methods("PUT")
	//  /quiz/{qid:[0-9]+}/splash PUT (i.e. teacher about to start quiz, accept connects, state -1?)
	//  /quiz/{qid:[0-9]+}/connect PUT (i.e. student subscribe w/ HOST:PORT)

	//Javascript pages
	r.HandleFunc("/dashboard/", handlePage("dashboard")).Methods("GET")
	//TODO browser "student" client

	http.Handle("/", r)
	serveFile("/favicon.ico", "./favicon.ico")
	http.ListenAndServe(":8080", nil)
}
