package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

var templates = template.Must(template.New("").Delims("<<<", ">>>").ParseGlob("templates/*.html"))
var db *sql.DB

//TODO figure out somewhere better to put this...
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

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.HandleFunc("/", handlePage("index"))
	r.HandleFunc("/login", handlePage("login")).Methods("GET")
	r.HandleFunc("/logmein", login).Methods("POST")
	r.HandleFunc("/logmeout", logout)
	r.HandleFunc("/register", handlePage("register")).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")

	//TODO these are just ideas

	r.HandleFunc("/quiz/list", handleQuizList).Methods("GET")
	r.HandleFunc("/quiz/update/{id}/{title}/{info}/{cid}", handleQuizUpdate).Methods("GET")
	r.HandleFunc("/quiz/{id}", handleQuizGet).Methods("GET")
	//r.HandleFunc("/quiz/{id}", handleAnswer).Methods("PUT")
	//r.HandleFunc("/quiz/{id}/edit", handleQuizEdit).Methods("GET")
	//r.HandleFunc("/quiz/add, handleQuizCreate).Methods("POST")

	r.HandleFunc("/dashboard", handlePage("dashboard")).Methods("GET")
	r.HandleFunc("/quiz", handlePage("quiz")).Methods("GET")
	r.HandleFunc("/quiz/add", handleQuizCreate).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
