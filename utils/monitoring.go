package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "code", "method"},
	)
)

func init() {
	prometheus.MustRegister(HTTPRequestTotal)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}