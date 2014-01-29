package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

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
	if email == "" || pass1 == "" || pass2 == "" || pass1 != pass2 {
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
		MaxAge:  -1,
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
func auth(w http.ResponseWriter, r *http.Request) *User {
	cookie, err := r.Cookie("logged-in")
	if err != nil { // no cookie
		return nil
	} else { //cookie
		var uid int
		err = db.QueryRow(`SELECT uid FROM users WHERE email=$1`, cookie.Value).Scan(&uid)
		if err != nil { // invalid user
			return nil
		}
		return &User{uid, cookie.Value} // valid user
	}
}
