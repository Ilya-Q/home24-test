package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ilya-Q/home24-test/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("cfg")

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 8080)
	viper.SetDefault("timeout", 5*time.Second)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/", handlers.AnalysisFormHandler).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.NotFoundHandler)

	err = http.ListenAndServe(
		fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port")),
		r,
	)
	if err != nil {
		log.Fatalf("Couldn't start HTTP server: %v", err)
	}
}
