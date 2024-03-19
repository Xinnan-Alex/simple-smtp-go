package main

import (
	"bytes"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/smtp"
	"strings"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type Mail struct {
	subject      string
	body         string
	templatePath string
	to           []string
	cc           []string
	bcc          []string
}

type EmailSender interface {
	SendEmail(
		mail Mail,
	) error
	ToBytes(
		mail Mail,
	) []byte
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func (g *GmailSender) ToBytes(mail Mail) []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf("From: %s\r\n", g.fromEmailAddress))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.to, ",")))
	if len(mail.cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.cc, ",")))
	}

	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.subject))

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n\n", boundary))
	buf.WriteString(fmt.Sprintf("--%s\n", boundary))

	buf.WriteString("Content-Type: text/html; charset=utf-8\n")

	buf.WriteString(mail.body)
	buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))

	return buf.Bytes()
}

func (g *GmailSender) SendEmail(mail Mail) error {
	// Get HTML

	if len(mail.templatePath) > 0 {
		var body bytes.Buffer
		t, err := template.ParseFiles(mail.templatePath)
		err = t.Execute(&body, struct {
			Name string
		}{Name: "Robby"})
		if err != nil {
			fmt.Println(err)
			return err
		}
		mail.body = body.String()
	}

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		g.fromEmailAddress,
		g.fromEmailPassword,
		smtpAuthAddress,
	)

	// Append all To,CC,BCC to the all slice
	var all []string
	for _, a := range [][]string{mail.to, mail.cc, mail.bcc} {
		all = append(all, a...)
	}

	emailMessageBytes := g.ToBytes(
		mail,
	)

	err := smtp.SendMail(smtpServerAddress,
		auth,
		g.fromEmailAddress,
		all,
		emailMessageBytes,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func NewGmail(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}
