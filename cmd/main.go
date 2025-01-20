package main

import (
	"SAIL-user-service/config"
	"SAIL-user-service/handlers"
	"SAIL-user-service/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	router := mux.NewRouter()

	handlers.RegisterUserHandlers(router, cfg)

	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	router.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})).Methods("GET")

	router.Use(middleware.PrometheusMiddleware)


	log.Printf("Starting server on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}