package main

import (
	"crypto/tls"

	log "github.com/sirupsen/logrus"

	gomail "gopkg.in/mail.v2"
)

type Mail struct {
	To      string
	Body    string
	Subject string
}

func NewMail(to, body, subject string) *Mail {
	return &Mail{
		To:      to,
		Body:    body,
		Subject: subject,
	}
}

type Mailer struct {
	host     string
	port     int
	username string
	password string
	mailFrom string
}

func NewMailer(host string, port int, username string, password string, mailfrom string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		mailFrom: mailfrom,
	}
}

func (ml *Mailer) SendMail(mail *Mail) (bool, error) {

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", ml.mailFrom)

	// Set E-Mail receivers
	m.SetHeader("To", mail.To)

	// Set E-Mail subject
	m.SetHeader("Subject", mail.Subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", mail.Body)

	// Settings for SMTP server
	d := gomail.NewDialer(ml.host, ml.port, ml.username, ml.password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		log.Errorln("Error while sending email:" + err.Error())
		return false, err
	}

	return true, nil
}
