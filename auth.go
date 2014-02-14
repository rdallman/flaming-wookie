package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/securecookie"
	_ "github.com/lib/pq"
)

var hashKey = []byte("notthesecretyourelookingfor")
var blockKey = []byte("chewbacachewbaca")
var s = securecookie.New(hashKey, blockKey)

// login gets the user' email and password from the form.
// If valid user, createCookie(w, inputEmail) is called, and
// the user is redirected to their dashboard page.
func login(w http.ResponseWriter, r *http.Request) {
	inputEmail, inputPass := r.FormValue("email"), r.FormValue("password")

	//dbpw is the salted sha256 hash we stored as password
	var salt, dbpw string
	var uid int
	err := db.QueryRow(`SELECT salt, password, uid  FROM users WHERE email=$1`, inputEmail).Scan(&salt, &dbpw, &uid)

	//salt input password, hash and compare to database salted hash
	hash := sha256.Sum256(append([]byte(inputPass), salt...))
	phash := string(hash[:]) //finicky

	switch {
	case err == sql.ErrNoRows, dbpw != phash:
		// TODO add flash messages
		http.Redirect(w, r, "/login", 302)
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		err = createCookie(w, uid, inputEmail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/dashboard/#/main", 302)
		}
	}
}

// register gets the user's email and password from the form.
// If the fields are valid and user does not already exist, the
// user is added to the database, createCookie(w, inputEmail) is
// called, and the user is redirected to their dashboard page.
func register(w http.ResponseWriter, r *http.Request) {
	email, pass1, pass2 := r.FormValue("email"), r.FormValue("password"), r.FormValue("password confirm")

	//check if any fields are blank or if passwords do not match
	if email == "" || pass1 == "" || pass2 == "" || pass1 != pass2 {
		fmt.Fprintf(w, "Invalid information")
		return //do not add to db
	}
	//check to see if email already exists in db
	var id int
	err := db.QueryRow(`SELECT uid FROM users WHERE email=$1`, email).Scan(&id)
	if id > 0 {
		//FIXME flash message and try again, this is bad
		fmt.Fprintf(w, "User already exists")
		return //do not add to db
	}

	salt := make([]byte, 32)
	_, err = rand.Read(salt)

	// create secure, salted hash
	hash := sha256.Sum256(append([]byte(pass1), salt...))
	phash := string(hash[:])

	var uid int
	err = db.QueryRow(`INSERT INTO users (email, password, salt)
    VALUES($1, $2, $3) RETURNING uid`, email, phash, salt).Scan(&uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//send to login
	http.Redirect(w, r, "/login", 302)
}

// createCookie creates and sets a cookie for the user.
func createCookie(w http.ResponseWriter, uid int, email string) error {
	values := make(map[string]string)
	values["uid"] = strconv.Itoa(uid)
	values["email"] = email
	encoded, err := s.Encode("logged-in", values)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:    "logged-in",
		Value:   encoded,
		Path:    "/",
		Expires: time.Now().AddDate(20, 0, 0),
	}

	http.SetCookie(w, cookie)
	return nil
}

// logout deletes cookie and redirects to homepage.
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("logged-in")
	if err != nil { //no cookie
		return
	}
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

// auth checks for a user's cookie. If a valid cookie exists,
// a User struct is returned containing the user's uid and email.
// If no valid cookie exists, nil is returned.
func auth(r *http.Request) *User {
	cookie, err := r.Cookie("logged-in")
	if err != nil { // no cookie
		return nil
	}
	values := make(map[string]string)
	err = s.Decode("logged-in", cookie.Value, &values)
	if err != nil { // invalid user
		return nil
	}

	id, err := strconv.Atoi(values["uid"])
	if err != nil {
		return nil
	}

	return &User{id, values["email"]} // valid user
}
