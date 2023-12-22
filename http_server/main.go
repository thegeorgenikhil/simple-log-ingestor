package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	RMQPublisherClient *RabbitMQPublisher
)

const port = "9119"

func main() {
	// connect to rabbitmq
	rcl, err := NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln("can't connect to rabbitmq", err)
	}

	RMQPublisherClient = rcl

	fmt.Println("connected to rabbitmq...")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: Routes(),
	}

	fmt.Println("Starting server at port", port)
	log.Fatalln(srv.ListenAndServe())
}
