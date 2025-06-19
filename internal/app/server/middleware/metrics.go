package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com.br/gibranct/simplified-wallet/internal/provider/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// PrometheusMiddleware collects metrics for HTTP requests
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Call the next handler
		next.ServeHTTP(ww, r)

		// Record metrics after the request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ww.Status())
		endpoint := chi.RouteContext(r.Context()).RoutePattern()
		if endpoint == "" {
			endpoint = r.URL.Path
		}

		metrics.RequestCounter.WithLabelValues(r.Method, endpoint, status).Inc()
		metrics.RequestDuration.WithLabelValues(r.Method, endpoint).Observe(duration)
	})
}
