package server

import (
	"log"
	"net/http"

	"github.com.br/gibranct/simplified-wallet/internal/app/server/router"
	"github.com.br/gibranct/simplified-wallet/internal/config"
	"github.com/golang-migrate/migrate/v4"
)

func Run() {
	dbConfig := config.GetPostgresConfig()
	m, err := migrate.New(
		"file://migrations",
		dbConfig.GetPostgresURL(),
	)
	if err != nil {
		log.Fatal("Failed to initialize migration, err: ", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to apply migrations, err: ", err)
	}
	r := router.InitRoutes()
	log.Println("Server running on port 3000...")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal("Failed to start server, err: ", err)
	}
}
