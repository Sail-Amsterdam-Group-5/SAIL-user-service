package main

import (
	"SAIL-user-service/config"
	"SAIL-user-service/handlers"
	"SAIL-user-service/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title SAIL User Service API
// @version 1.0
// @description This is the API for the SAIL User Service
// @termsOfService https://swagger.io/terms/
// @contact.name API Support
// @contact.url https://www.sail.com
// @contact.email email.placeholder@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
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

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler) //since it's system wide docs

	router.Use(middleware.PrometheusMiddleware)


	log.Printf("Starting server on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}