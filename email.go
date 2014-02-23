package main

import (
	"net/smtp"
	"bytes"
	"text/template"
	"log"
	"strconv"
	"encoding/json"

	_ "github.com/lib/pq"
)

type EmailUser struct {
  Username    string
  Password    string
  EmailServer string
  Port        int
}

var emailUser = &EmailUser{"wookiet3st", "wooquiz1", "smtp.gmail.com", 587}

var emailauth = smtp.PlainAuth("",
  emailUser.Username,
  emailUser.Password,
  emailUser.EmailServer,
)

type SmtpTemplateData struct {
  From	  string
  Subject string
  Body    string
}

const emailTemplate = `From: {{.From}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,
{{.From}}
`

func sendStudentEmail(cid int) {
	var err error
	
	//get stuff from db
	var classname, studentjson string
	err = db.QueryRow(`SELECT name, students FROM classes WHERE cid=$1`, cid).Scan(&classname, &studentjson)
	if err != nil {
		return
	}
	var students []map[string]string
	json.Unmarshal([]byte(studentjson), &students)

	//loop over students and send email
	for _, student := range students {
		var doc bytes.Buffer
		context := SmtpTemplateData{
		  "WooQuiz",
		  classname + " Class Registration",
		  "Hello " + student["name"] + ", You have been added to this class. Here is your awesome stuff that you can do stuff with! yay",
		}

		t := template.New("emailTemplate")
		t, err = t.Parse(emailTemplate)
		if err != nil {
		  log.Print("error trying to parse mail template")
		}
		err = t.Execute(&doc, context)
		if err != nil {
		  log.Print("error trying to execute mail template")
		}

		err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
		  emailauth,
		  emailUser.Username,
		  []string{student["email"]},
		  doc.Bytes())
		if err != nil {
		  log.Print("ERROR: attempting to send a mail ", err)
		}
	}
}