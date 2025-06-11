package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("simplified-wallet")

// TracingMiddleware adds OpenTelemetry tracing to HTTP requests
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract tracing context from the request headers
		ctx := r.Context()
		propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

		// Get the route pattern
		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		if routePattern == "" {
			routePattern = r.URL.Path
		}

		// Start a new span
		opts := []trace.SpanStartOption{
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.route", routePattern),
			),
		}
		ctx, span := tracer.Start(ctx, "HTTP "+r.Method+" "+routePattern, opts...)
		defer span.End()

		// Update the request with the new context
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
