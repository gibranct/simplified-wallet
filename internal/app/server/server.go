package server

import (
	"context"
	"errors"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"log"
	"net/http"

	"github.com.br/gibranct/simplified-wallet/internal/app/server/router"
	"github.com.br/gibranct/simplified-wallet/internal/config"
	"github.com/golang-migrate/migrate/v4"
)

const serviceName = "simplified-wallet"

func Run() {
	// Initialize tracing
	otel, err := telemetry.NewJaeger(context.Background(), serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := otel.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
		log.Println("OTLP tracer provider shut down")
	}()

	dbConfig := config.GetPostgresConfig()
	m, err := migrate.New(
		"file://migrations",
		dbConfig.GetPostgresURL(),
	)
	if err != nil {
		log.Fatal("Failed to initialize migration, err: ", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Failed to apply migrations, err: ", err)
	}
	r := router.InitRoutes(otel)
	log.Println("Server running on port 3000...")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal("Failed to start server, err: ", err)
	}
}
