package router

import (
	"net/http"
	"time"

	"github.com.br/gibranct/simplified-wallet/internal/app/server/handler"
	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/provider/db"
	"github.com.br/gibranct/simplified-wallet/internal/provider/gateway"
	repository "github.com.br/gibranct/simplified-wallet/internal/provider/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	createTransaction := usecase.NewCreateTransaction(
		repository.NewUserRepository(db.NewPostgresDB()),
		gateway.NewTransactionAuthorizer(http.DefaultClient),
	)

	h := handler.New(&createTransaction)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/transaction", h.PostTransaction)
	})
	return r
}
