package config

import (
	"fmt"
	"strings"
)

type PostgresConfig struct {
	URI     string
	SSLMode bool
}

func GetPostgresConfig() PostgresConfig {
	return PostgresConfig{
		URI:     getEnv("POSTGRES_URI", "postgres://postgres:postgres@localhost:5432/wallet"),
		SSLMode: getEnvAsBool("POSTGRES_SSLMODE", false),
	}
}

func (p PostgresConfig) GetPostgresURL() string {
	sslMode := "disable"
	if p.SSLMode {
		sslMode = "require"
	}
	if strings.HasSuffix(p.URI, "?sslmode=disable") {
		return p.URI
	}
	return fmt.Sprintf(
		"%s?sslmode=%s",
		p.URI,
		sslMode,
	)
}
