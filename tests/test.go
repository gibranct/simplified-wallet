package test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	testcontainers.Container
	URI string
}

func setupPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "wallet_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("postgres://postgres:postgres@%s:%s/wallet_test?sslmode=disable", hostIP, mappedPort.Port())

	return &PostgresContainer{
		Container: container,
		URI:       uri,
	}, nil
}

func runMigrations(db *sqlx.DB, version uint) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func SetupTestDatabase(ctx context.Context, migrateVersion uint) (*PostgresContainer, *sqlx.DB, error) {
	container, err := setupPostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Wait a bit for the container to be fully ready
	time.Sleep(2 * time.Second)

	db, err := sqlx.Connect("postgres", container.URI)
	if err != nil {
		return container, nil, err
	}
	err = os.Setenv("POSTGRES_URI", container.URI)
	if err != nil {
		return container, db, err
	}

	if err := db.Ping(); err != nil {
		return container, db, err
	}

	if err := runMigrations(db, migrateVersion); err != nil {
		return container, db, err
	}

	return container, db, nil
}
