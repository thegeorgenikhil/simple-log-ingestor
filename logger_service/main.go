package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// The goal with log_consumer is to consume logs from RabbitMQ and write them to a file.

const (
	_logFile      = "all_logs.log"
	_exchangeName = "logs_exchange"
)

func main() {
	logger := NewLogger(_logFile)

	routingKeys := []string{
		"info",
		"debug",
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

			switch l.Level {
			case "info":
				logger.Info(l.Message, l.From, l.Time)
			case "debug":
				logger.Debug(l.Message, l.From, l.Time)
			case "error":
				logger.Error(l.Message, l.From, l.Time)
			default:
				fmt.Println("[LOG_CONSUMER] unknown log level: ", l.Level)
			}

			d.Ack(false)
		}
	}()

	<-forever
}
