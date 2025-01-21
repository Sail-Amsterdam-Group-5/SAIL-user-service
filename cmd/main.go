package main

import (
	"SAIL-user-service/config"
	"SAIL-user-service/handlers"
	"SAIL-user-service/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{}) // JSON format for structured logs
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Starting SAIL User Service API")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Error loading config")
	}

	logger.WithField("port", cfg.Port).Info("Configuration loaded")

	router := mux.NewRouter()

	handlers.RegisterUserHandlers(router, cfg)
	logger.Info("Registered user handlers")

	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	logger.Info("Registered /metrics endpoint")

	router.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})).Methods("GET")
	logger.Info("Registered /health endpoint")

	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	logger.Info("Serving static files from /docs")

	router.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"),
	)).Methods("GET")
	logger.Info("Registered Swagger documentation endpoint")


	router.Use(middleware.PrometheusMiddleware)
	logger.Info("Middleware applied")


	logger.WithField("port", cfg.Port).Info("Starting server")
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.WithError(err).Fatal("Server failed to start")
	}
}