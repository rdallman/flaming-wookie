package main

import (
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
