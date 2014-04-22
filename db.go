package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// init opens a connection to the database.
func init() {
	var err error
	//db, err = sql.Open("postgres", "postgres://wookie:password@absker.com/wookie?sslmode=disable")
	db, err = sql.Open("postgres", "user=wookie dbname=wookie sslmode=disable")
	if err != nil {
		ERROR.Println("init db", err.Error())
		return
	}

	//////////////////////////////////////
	// drop tables
	// DANGER this will empty the db
	//
	//////////////////////////////////////
	_, err = db.Exec(`DROP TABLE classes, users, quiz, attendance CASCADE`)
	fmt.Println(err)

	/////////////////////////////////////////////
	//////creating
	/////////////////////////////////////////////

	_, err = db.Exec(`CREATE TABLE users (
    uid serial PRIMARY KEY,
    email text UNIQUE,
    password bytea,
    salt bytea
  )`)
	fmt.Println(err)

	_, err = db.Exec(`CREATE TABLE attendance (
    cid           integer PRIMARY KEY,
    students      json,
    date_created  date
  )`)
	fmt.Println(err)

	_, err = db.Exec(`CREATE TABLE classes (
    cid serial PRIMARY KEY,
    name text,
    students json,
    uid integer REFERENCES users (uid),
    semester text
  )`)
	fmt.Println(err)

	_, err = db.Exec(`CREATE TABLE quiz (
  qid serial PRIMARY KEY,
  info json,
  type integer,
  cid integer REFERENCES classes (cid)
  )`)
	fmt.Println(err)
}
