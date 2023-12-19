package util

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func InitDialer() {
	host := os.Getenv("MAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		Logger.Error(err)
	}
	userName := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	dialer = gomail.NewDialer(host, port, userName, password)
}

func SendMail(message *gomail.Message) error {
	message.SetHeader("From", dialer.Username)
	return dialer.DialAndSend(message)
}
