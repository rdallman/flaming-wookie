package main

import (
	"bytes"
	"fmt"
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

	cidStr := strconv.Itoa(cid)
	link := fmt.Sprintf(`http://%s:%s@wooquiz.com`, student["sid"], cidStr)
	//create text for email
	var doc bytes.Buffer
	context := SmtpTemplateData{
		From:    "WooQuiz",
		Subject: fmt.Sprintf("%s Class Registration", classname),
		Body: fmt.Sprintf(`Hello %s %s,

You have been added to the class %s on WooQuiz.com!

  Class ID: %s 
  Student ID: %s

Either enter the above information or click the link below on your mobile to automatically add the class.

%s

Good Luck!

--
WooQuiz`, student["fname"], student["lname"], classname, cidStr, student["sid"], link),
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
