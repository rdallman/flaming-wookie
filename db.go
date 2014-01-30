package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// init opens a connection to the database.
func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://wookie:password@absker.com/wookie?sslmode=disable")
	//db, err = sql.Open("postgres", "user=reed dbname=wookie sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
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

	////for getting table names -- handy

	//rows, err := db.Query("SELECT tablename from pg_catalog.pg_tables")
	//for rows.Next() {
	//var tablename string
	//rows.Scan(&tablename)
	//fmt.Println(tablename)
	//}

	/////////////////////////////////////////////
	//////creating
	//////TODO NOT NULL all of these later...
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
