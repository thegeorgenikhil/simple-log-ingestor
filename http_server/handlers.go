package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type LogRequest struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	From    string `json:"from"`
}

// IngestLog takes the received log in the request body and pushes it to the RabbitMQ
func IngestLog(w http.ResponseWriter, r *http.Request) {
	var lr LogRequest
	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		fmt.Println("Error while decoding request body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body not valid."))
		return
	}

	if lr.Level == "" || lr.Message == "" || lr.From == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body not valid."))
		return
	}

	formattedLogLevel := strings.ToLower(lr.Level)
	if formattedLogLevel != "info" && formattedLogLevel != "debug" && formattedLogLevel != "error" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("log level should be one of info, debug or error."))
		return
	}

	log := &LogData{
		Level:   formattedLogLevel,
		Message: lr.Message,
		From:    lr.From,
		Time:    time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = RMQPublisherClient.publishLog(ctx, formattedLogLevel, log)

	if err != nil {
		fmt.Println("Error while publishing log: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while publishing log."))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Log received!"))
}
