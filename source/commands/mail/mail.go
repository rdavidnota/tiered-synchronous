package mail

import (
	"log"
	"net/smtp"
	"strings"
)

func SendMail(from string, to []string, message string) {

	message = "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + mailConfig["subject"] + "\n\n" + message

	err := smtp.SendMail(mailConfig["host"]+":"+mailConfig["port"],
		smtp.PlainAuth("", mailConfig["username"], mailConfig["password"], mailConfig["host"]),
		mailConfig["username"], to, []byte(message))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, From: " + from + ", To:" + strings.Join(to, ","))
}
