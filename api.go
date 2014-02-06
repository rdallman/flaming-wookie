package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

// add more later
func handleQuizGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Printf("%d", qID) //testing
		rows, err := db.Query(`SELECT * FROM quiz WHERE qid=$1`, qID)
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			for rows.Next() {
				//fmt.Printf("here") //testing
				var qid, cid int
				var title, info string
				err = rows.Scan(&qid, &title, &info, &cid)
				if err != nil {
					fmt.Printf("%s", err)
				} else {
					fmt.Printf("\nqid:%d \ttitle:%s \tinfo:%s \tcid:%d", qid, title, info, cid)
				}
			}
		}
	}
}

// just an idea, not sure if we actually need this
func handleQuizUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HERE")
	vars := mux.Vars(r)
	//not sure we need all these, but for now...
	qid, err1 := strconv.Atoi(vars["id"])
	title := r.FormValue("title")
	info := r.FormValue("info")
	cid, err2 := strconv.Atoi(r.FormValue("cid"))

	if err1 != nil || err2 != nil {
		fmt.Println("err1 $1\terr2 $2", err1, err2)
	} else {
		_, err := db.Exec(`UPDATE quiz SET title=$1, info=$2, cid=$3 WHERE qid=$4`, title, info, cid, qid)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("\nUpdated.")
		}
	}

}

func handleQuizDelete(w http.ResponseWriter, r *http.Request) {
  authenticate := auth(r)
  if authentication != nil { 
  vars := mux.Vars(r)
  qid, err := strconv.Atoi(vars["id"])
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(qid)
    _, err := db.Exec(`DELETE FROM quiz WHERE qid=$1`, qid)
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("\nDeleted.")
    }
  }
  } else { //this is bad, but we can decide this later...
    fmt.Println("Error - you cannot delete a quiz")
  }
}

func handleQuizList(w http.ResponseWriter, r *http.Request) {
	//title and id, return JSON
	rows, err := db.Query(`SELECT qid, title FROM quiz`)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		quizzes := make(map[string]int)
		for rows.Next() {
			var qid int
			var title string
			err = rows.Scan(&qid, &title)
			if err != nil {
				fmt.Println(err)
			} else {
				quizzes[title] = qid
			}
		}
		jquiz, err := json.Marshal(quizzes)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Fprintf(w, string(jquiz))
		}
	}
}
