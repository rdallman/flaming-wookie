package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

//checks for err, replies with false success reply
func writeErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return true
	}
	return false
}

//TODO want support for multiple items?
func writeSuccess(w http.ResponseWriter, info ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	r := Response{"success": true}
	if len(info) == 1 {
		r["info"] = info[0]
	}
	fmt.Fprint(w, r)
}

type Response map[string]interface{}

func (r Response) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}
