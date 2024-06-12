package main

import (
	"net/http"

	"github.com/Ilya-Q/home24-test/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/", handlers.AnalysisFormHandler).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.NotFoundHandler)

	http.ListenAndServe("localhost:8080", r)
}
