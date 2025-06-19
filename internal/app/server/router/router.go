package router

import (
	"github.com.br/gibranct/simplified-wallet/internal/provider/queue"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"

	"github.com.br/gibranct/simplified-wallet/internal/app/server/handler"
	customMiddleware "github.com.br/gibranct/simplified-wallet/internal/app/server/middleware"
	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/app/usecase/strategy"
	"github.com.br/gibranct/simplified-wallet/internal/provider/db"
	"github.com.br/gibranct/simplified-wallet/internal/provider/gateway"
	repository "github.com.br/gibranct/simplified-wallet/internal/provider/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRoutes(otel telemetry.Telemetry) *chi.Mux {
	r := chi.NewRouter()

	// Standard middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Add custom middleware for metrics and tracing
	r.Use(customMiddleware.PrometheusMiddleware)
	r.Use(customMiddleware.TracingMiddleware)

	// Expose Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	userRepo := repository.NewUserRepository(db.NewPostgresDB(), otel)
	createTransaction := usecase.NewCreateTransaction(
		userRepo,
		gateway.NewTransactionAuthorizer(http.DefaultClient, otel),
		queue.NewSNS(otel),
		otel,
	)
	strategies := []usecase.CreateUserStrategy{
		strategy.NewCreateCommonUser(userRepo, otel),
		strategy.NewCreateMerchantUser(userRepo, otel),
	}
	createUser := usecase.NewCreateUser(userRepo, strategies, otel)

	h := handler.New(createTransaction, createUser, otel)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/transactions", h.PostTransaction)
		r.Post("/users", h.PostUser)
		r.Post("/merchants", h.PostMerchant)
	})
	return r
}
