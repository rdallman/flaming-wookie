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
	r.HandleFunc("/login", handlePage("login")).Methods("GET")
	r.HandleFunc("/logmein", login).Methods("POST")
	r.HandleFunc("/logmeout", logout)
	r.HandleFunc("/register", handlePage("register")).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")

	//TODO these are just ideas

	r.HandleFunc("/dashboard/quiz", handleQuizList).Methods("GET")
	r.HandleFunc("/dashboard/quiz", handleQuizCreate).Methods("POST")
	r.HandleFunc("/dashboard/quiz/{id}", handleQuizUpdate).Methods("POST")
	r.HandleFunc("/dashboard/quiz/{id}", handleQuizGet).Methods("GET")
	//r.HandleFunc("/dashboard/quiz/{id}", handleQuizDelete).Methods("DELETE")

	r.HandleFunc("/quiz/{id}/state", changeState).Methods("PUT")
	r.HandleFunc("/quiz/{id}/answer", handleAnswer).Methods("PUT")

	r.HandleFunc("/dashboard/", handlePage("dashboard")).Methods("GET")
	r.HandleFunc("/quiz", handlePage("quiz")).Methods("GET")
	r.HandleFunc("/quiz/add", handleQuizCreate).Methods("POST")

	http.Handle("/", r)
	serveFile("/favicon.ico", "./favicon.ico")
	http.ListenAndServe(":8080", nil)
}
