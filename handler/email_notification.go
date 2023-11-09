package handler

import (
	"bytes"
	"fmt"
	"hotel/dto"
	"html/template"

	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(data dto.MailData) {
	// get html
	var body bytes.Buffer
	t, err := template.ParseFiles("./template/mail.html")
	t.Execute(&body, data)

	if err != nil {
		fmt.Println(err)
	}

	// Set Header
	m := gomail.NewMessage()
	m.SetHeader("From", "jacksparrow257257@gmail.com")
	m.SetHeader("To", "osvaldosilitonga@gmail.com")
	m.SetHeader("Subject", "Hotel Mini Project")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, "jacksparrow257257@gmail.com", os.Getenv("EMAIL_PASSWORD"))

	// Send
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
