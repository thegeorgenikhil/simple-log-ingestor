package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/log", IngestLog)

	return r
}
