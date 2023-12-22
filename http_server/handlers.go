package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LogRequest struct {
	Level string `json:"level"`
	Message string `json:"message"`
	From string `json:"from"`
}

// IngestLog takes the received log in the request body and pushes it to the RabbitMQ
func IngestLog(w http.ResponseWriter, r *http.Request) {
	var lr LogRequest
	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		fmt.Println("Error while decoding request body: ",err)
		w.WriteHeader(http.StatusBadRequest)	
		w.Write([]byte("request body not valid."))
		return 
	}
	
	res, _ := json.MarshalIndent(lr, "", "  ")
	fmt.Println(string(res))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Log received!"))
}