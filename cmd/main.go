package main

import (
	"github.com.br/gibranct/simplified-wallet/internal/app/server"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	server.Run()
}
