package main

import (
	"bytes"
	"net/smtp"
	"strconv"
	"text/template"

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
	From    string
	Subject string
	Body    string
}

const emailTemplate = `From: {{.From}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,
{{.From}}
`

func sendStudentClassEmail(cid int, classname string, student map[string]string) {
	var err error

	//create text for email
	var doc bytes.Buffer
	context := SmtpTemplateData{
		"WooQuiz",
		classname + " Class Registration",
		"Hello " + student["fname"] + " " + student["lname"] + "," +
			"\n\nYou have been added to the class " + classname + " on WooQuiz.com!" +
			"\nClass ID: " + strconv.Itoa(cid) +
			"\nStudent ID: " + student["sid"],
	}

	//template email
	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		ERROR.Println("Send Student Email - error trying to parse mail template")
	}
	err = t.Execute(&doc, context)
	if err != nil {
		ERROR.Println("Send Student Email - error trying to execute mail template")
	}

	//send email
	err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
		emailauth,
		emailUser.Username,
		[]string{student["email"]},
		doc.Bytes())
	if err != nil {
		ERROR.Println("Send Student Email - attempting to send email", err.Error())
	}
}
