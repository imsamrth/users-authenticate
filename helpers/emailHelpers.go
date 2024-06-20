package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"instix_auth/constants"
	"net/smtp"
	"os"
)

var (
	fromEmail    = os.Getenv("FromEmail")
	SMTPpassword = os.Getenv("SMTPpassword")
	Domain       = os.Getenv("DOMAIN")
)

type Request struct {
	body string
}

func NewRequest(body string) *Request {
	return &Request{
		body: body,
	}
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func SendEmail(name, email, subject, HTMLbody, verficationLink string) error {

	//NOTE : CHANGE THE HARD CODED EMAIL BEFRORE COMMIT
	// sender data
	//email = "samarthbnsl@gmail.com"
	to := []string{email}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp-auth.iitb.ac.in"
	port := "587"

	address := host + ":" + port
	// Set up authentication information.
	auth := smtp.PlainAuth("", fromEmail, SMTPpassword, host)

	templateData := struct {
		Name string
		URL  string
		LOGO string
	}{
		Name: name,
		URL:  verficationLink,
		LOGO: Domain + constants.AssetsDir + "/logo.jpg",
	}

	r := NewRequest("Hello, World!")
	err := r.ParseTemplate("templates/private/emailVerification.html", templateData)

	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	msg := []byte(
		"From: " + constants.SENDER + ": <" + fromEmail + ">\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			r.body)
	fmt.Println(fromEmail, SMTPpassword)
	err = smtp.SendMail(address, auth, fromEmail, to, msg)
	fmt.Println(err)
	if err != nil {
		return err
	}
	fmt.Println("Check for sent email!")
	return nil
}
