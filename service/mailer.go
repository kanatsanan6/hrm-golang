package service

import (
	"log"
	"os"
	"strings"

	"github.com/kanatsanan6/hrm/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct{}

func removePlus(email string) string {
	if strings.Contains(email, "+") {
		atIdx := strings.Index(email, "@")
		return email[:strings.Index(email, "+")] + email[atIdx:]
	} else {
		return email
	}
}

func (m *Mailer) Send(to string, subject string, message *gomail.Message) {
	message.SetHeader("From", os.Getenv("MAILER_USERNAME"))
	message.SetHeader("To", removePlus(to))
	message.SetHeader("Subject", subject)

	if err := config.Mailer.DialAndSend(message); err != nil {
		log.Panicln("[Mailer] ", err)
	}
}
