package util

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func InitDialer() {
	host := os.Getenv("MAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("MAIL_HOST"))
	if err != nil {
		Logger.Error(err)
	}
	userName := os.Getenv("MAIL_HOST")
	password := os.Getenv("MAIL_HOST")
	dialer = gomail.NewDialer(host, port, userName, password)
}

func SendMail(message *gomail.Message) error {
	message.SetHeader("From", dialer.Username)
	return dialer.DialAndSend(message)
}
