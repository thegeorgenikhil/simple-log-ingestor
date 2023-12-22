package main

import amqp "github.com/rabbitmq/amqp091-go"

type RabbitMQConsumer struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	logChannel <-chan amqp.Delivery
}

func NewRabbitMQConsumer(connectionString string, exchangeName string, routingKeys []string) (*RabbitMQConsumer, error) {
	c := &RabbitMQConsumer{}
	var err error

	c.conn, err = amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = c.configureExchange(exchangeName)
	if err != nil {
		return nil, err
	}

	q, err := c.configureQueue()
	if err != nil {
		return nil, err
	}

	// TODO: QoS - Quality of Service: https://www.rabbitmq.com/consumer-prefetch.html

	for _, v := range routingKeys {
		err = c.bindQueue(q.Name, v, exchangeName)
		if err != nil {
			return nil, err
		}
	}

	err = c.startConsumer(q.Name)
	if err != nil {
		return nil, err
	}

	return c, err
}

func (c *RabbitMQConsumer) configureExchange(exchangeName string) error {
	err := c.ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	return err
}

func (c *RabbitMQConsumer) configureQueue() (amqp.Queue, error) {
	q, err := c.ch.QueueDeclare(
		"",    // name
		true, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return q, err
}

func (c *RabbitMQConsumer) bindQueue(queueName string, routingKey string, exchangeName string) error {
	err := c.ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	return err
}

func (c *RabbitMQConsumer) startConsumer(queueName string) error {
	lCh, err := c.ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	c.logChannel = lCh

	return err
}

func (c *RabbitMQConsumer) GetChannel() <-chan amqp.Delivery {
	return c.logChannel
}
