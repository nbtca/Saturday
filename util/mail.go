package util

import (
	"strconv"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func InitDialer() {
	host := viper.GetString("mail.host")
	port, err := strconv.Atoi(viper.GetString("mail.port"))
	if err != nil {
		Logger.Error(err)
	}
	userName := viper.GetString("mail.username")
	password := viper.GetString("mail.password")
	dialer = gomail.NewDialer(host, port, userName, password)
}

func SendMail(message *gomail.Message) error {
	message.SetHeader("From", dialer.Username)
	return dialer.DialAndSend(message)
}
