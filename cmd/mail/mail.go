package mail

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
)

func SendText(username string, password string, from string, to []string, subject string, body string) {
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		"smtp.gmail.com",
	)

	msg := "Subject: " + subject + "\n" + body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		to,
		[]byte(msg),
	)
	if err != nil {
		log.Printf("smtp error: %v", err)
		return
	}
}

func SendHTML(username string, password string, from string, to []string, name string, subject string, templatePath string) {
	var body bytes.Buffer
	t, parse_err := template.ParseFiles(templatePath)
	if parse_err != nil {
		log.Printf("parse_error: %v", parse_err)
		return
	}

	execute_error := t.Execute(&body, struct{ Name string }{Name: name})
	if execute_error != nil {
		log.Printf("execute_error: %v", execute_error)
		return
	}

	auth := smtp.PlainAuth(
		"",
		username,
		password,
		"smtp.gmail.com",
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	msg := "Subject: " + subject + "\n" + headers + "\n\n" + body.String()

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		to,
		[]byte(msg),
	)
	if err != nil {
		log.Printf("smtp error: %v", err)
		return
	}
}
