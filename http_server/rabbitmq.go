package main

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_exchangeName = "logs_exchange"
)

type RabbitMQPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQPublisher(connectionString string) (*RabbitMQPublisher, error) {
	c := &RabbitMQPublisher{}
	var err error

	c.conn, err = amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = c.configureExchange()

	return c, err
}

func (c *RabbitMQPublisher) configureExchange() error {
	err := c.ch.ExchangeDeclare(
		_exchangeName, // name
		"direct",      // type 
		// The routing algorithm behind a "direct" exchange is simple - a message goes to the queues whose routing key exactly matches the routing key of the message.
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	return err
}

func (c *RabbitMQPublisher) publishLog(ctx context.Context,routingKey string, body *LogData) error {
	b,_ := json.Marshal(body)

	err := c.ch.PublishWithContext(ctx,
		_exchangeName, // name of the exchange to publish
		routingKey, 
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: []byte(b),
		},
	)

	return err
}
