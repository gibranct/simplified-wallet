package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestCounter counts the number of HTTP requests
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// RequestDuration measures the duration of HTTP requests
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// TransactionCounter counts the number of transactions
	TransactionCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "transactions_total",
			Help: "Total number of transactions",
		},
	)

	// TransactionAmount tracks the total amount of transactions
	TransactionAmount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "transactions_amount_total",
			Help: "Total amount of transactions",
		},
	)
)
