package main

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

func sendAnEMail(subject string, message string) {
	log := logging.MustGetLogger("log")

	host := viper.GetString("email.smtp")
	login := viper.GetString("email.login")
	password := viper.GetString("email.password")
	hostport := host + ":" + viper.GetString("email.port")
	from := viper.GetString("email.from")
	to := viper.GetStringSlice("email.sendTo")
	now := time.Now()

	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = fmt.Sprintf(subject)
	e.Text = []byte(fmt.Sprintf("À %v, %s", now, message))
	if err := e.Send(hostport, smtp.PlainAuth("", login, password, host)); err != nil {
		log.Warningf("Unable to send an email to \"%s\": %v", err)
	} else {
		log.Debug("Email was sent")
	}
}
