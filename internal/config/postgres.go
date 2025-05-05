package config

import "fmt"

type PostgresConfig struct {
	URI     string
	SSLMode bool
}

func GetPostgresConfig() PostgresConfig {
	return PostgresConfig{
		URI:     getEnv("POSTGRES_URI", "postgres://localhost:postgres@postgres:5432/wallet?"),
		SSLMode: getEnvAsBool("POSTGRES_SSLMODE", false),
	}
}

func (p PostgresConfig) GetPostgresURL() string {
	sslMode := "disable"
	if p.SSLMode {
		sslMode = "require"
	}
	return fmt.Sprintf(
		"%s?sslmode=%s",
		p.URI,
		sslMode,
	)
}
