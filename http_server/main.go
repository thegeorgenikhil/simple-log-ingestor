package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "9119"

func main(){
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",port),
		Handler: Routes(),
	}

	fmt.Println("Starting server at port",port)
	log.Fatal(srv.ListenAndServe())
}