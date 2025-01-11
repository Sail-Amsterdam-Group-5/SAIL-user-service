package main

import (
	"SAIL-user-service/config"
	"SAIL-user-service/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: %v", err)
	}

	router := mux.NewRouter()

	handlers.RegisterUserHandlers(router, cfg)

	log.Printf("Starting server on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}