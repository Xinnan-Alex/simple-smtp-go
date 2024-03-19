package main

import (
	"log"
)

func main() {

	// Provide Google SMTP Credentials
	identiy := ""
	fromEmailAddress := ""
	fromEmailPassword := ""

	mail := Mail{
		subject:      "Test Subject",
		templatePath: "./test.html",
		//body: "Test Body",
		to:  []string{""},
		cc:  []string{""},
		bcc: []string{""},
	}

	gmailSMTP := NewGmail(
		identiy,
		fromEmailAddress,
		fromEmailPassword,
	)

	err := gmailSMTP.SendEmail(mail)
	if err != nil {
		log.Panic(err)
		return
	}
}
