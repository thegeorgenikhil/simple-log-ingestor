package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// The goal with log_consumer is to consume logs from RabbitMQ and write them to a file.

const (
	_logFile      = "mails_sent.log"
	_exchangeName = "logs_exchange"
)

func main() {
	godotenv.Load()

	host := os.Getenv("EMAIL_SERVER_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_SERVER_PORT"))
	username := os.Getenv("EMAIL_SERVER_USERNAME")
	password := os.Getenv("EMAIL_SERVER_PASSWORD")
	mailFrom := os.Getenv("EMAIL_FROM_ADDRESS")
	mailTo := os.Getenv("EMAIL_TO_ADDRESS")

	logger := NewLogger(_logFile)
	mailer := NewMailer(
		host,
		port,
		username,
		password,
		mailFrom,
	)

	routingKeys := []string{
		"error",
	}

	RMQConsumerClient, err := NewRabbitMQConsumer("amqp://guest:guest@localhost:5672/", _exchangeName, routingKeys)
	if err != nil {
		log.Fatalln("[LOG_CONSUMER]can't connect to rabbitmq", err)
	}

	fmt.Println("[LOG_CONSUMER] connected to rabbitmq...")

	var forever chan struct{}

	logs := RMQConsumerClient.GetChannel()

	fmt.Println("[LOG_CONSUMER] started consuming logs from the queue...")

	go func() {
		for d := range logs {
			var l LogData
			_ = json.Unmarshal(d.Body, &l)

			content, err := ParseTemplate(ErrorReportingTemplateFile, l)
			if err != nil {
				fmt.Println("[LOG_CONSUMER] error parsing template", err)
				continue
			}

			mail := NewMail(mailTo, content, ErrorReportingSubject)

			ok, err := mailer.SendMail(mail)
			if err != nil {
				fmt.Println("[LOG_CONSUMER] error sending mail", err)
				continue
			}

			if ok {
				logger.Error(l.Message, l.From, l.Time, mailTo, time.Now())
			}

			d.Ack(false)
		}
	}()

	<-forever
}
