package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	repository "github.com.br/gibranct/simplified-wallet/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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

func setupTestDatabase(ctx context.Context, migrateVersion uint) (*PostgresContainer, *sqlx.DB, error) {
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

	if err := db.Ping(); err != nil {
		return container, db, err
	}

	if err := runMigrations(db, migrateVersion); err != nil {
		return container, db, err
	}

	return container, db, nil
}

func createTestUser(ctx context.Context, db *sqlx.DB, name, userType, document string, initialBalance float64) (uuid.UUID, error) {
	userID := uuid.New()
	money, err := vo.NewMoney(initialBalance)
	if err != nil {
		log.Panicln(err)
	}

	var query string

	if len(document) == 11 {
		query = `INSERT INTO users (id, name, cpf, email, password, user_type, balance) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	} else if len(document) == 14 {
		query = `INSERT INTO users (id, name, cnpj, email, password, user_type, balance) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	}

	_, err = db.ExecContext(
		ctx,
		query,
		userID,
		name,
		document,
		fmt.Sprintf("%s@example.com", name),
		"password",
		userType,
		money.Value(),
	)

	return userID, err
}

func getBalance(ctx context.Context, db *sqlx.DB, userID uuid.UUID) (int64, error) {
	var balance int64
	err := db.QueryRowContext(ctx, "SELECT balance FROM users WHERE id = $1", userID.String()).Scan(&balance)
	return balance, err
}

func TestCreateTransaction_Integration_Success(t *testing.T) {
	ctx := context.Background()
	migrateVersion := uint(2) // Use the latest migration version

	// Setup
	container, db, err := setupTestDatabase(ctx, migrateVersion)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer db.Close()
	cpf := "86395839004"
	cnpj := "71627571000107"

	// Create test users
	senderID, err := createTestUser(ctx, db, "sender", "common", cpf, 1000.0)
	require.NoError(t, err)

	receiverID, err := createTestUser(ctx, db, "receiver", "merchant", cnpj, 0.0)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(db)
	authorizerGateway := NewMockTransactionAuthorizerGateway(true) // Always authorize

	// Create use case
	createTransactionUseCase := usecase.NewCreateTransaction(userRepo, authorizerGateway)

	// Execute transaction
	amount := 400.0
	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	transactionID, err := createTransactionUseCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, transactionID)

	// Verify balances were updated correctly
	senderBalance, err := getBalance(ctx, db, senderID)
	require.NoError(t, err)
	assert.Equal(t, int64(60000), senderBalance)

	receiverBalance, err := getBalance(ctx, db, receiverID)
	require.NoError(t, err)
	assert.Equal(t, int64(40000), receiverBalance)

	// Verify transaction was recorded
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions WHERE id = $1", transactionID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCreateTransaction_Integration_Rollback(t *testing.T) {
	ctx := context.Background()
	usersTableVersion := uint(1) // Version of the users table where the error occurs

	// Setup
	container, db, err := setupTestDatabase(ctx, usersTableVersion)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer db.Close()
	cpf := "86395839004"
	cnpj := "71627571000107"

	// Create test users
	senderID, err := createTestUser(ctx, db, "sender_rollback", "common", cpf, 100.0)
	require.NoError(t, err)

	receiverID, err := createTestUser(ctx, db, "receiver_rollback", "merchant", cnpj, 0.0)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(db)

	authorizerGateway := NewMockTransactionAuthorizerGateway(true) // Always authorize

	// Create use case with the failing repository
	createTransactionUseCase := usecase.NewCreateTransaction(userRepo, authorizerGateway)

	// Get initial balances
	initialSenderBalance, err := getBalance(ctx, db, senderID)
	require.NoError(t, err)

	initialReceiverBalance, err := getBalance(ctx, db, receiverID)
	require.NoError(t, err)

	// Execute transaction that should fail
	amount := 50.0
	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	transactionID, err := createTransactionUseCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, transactionID)

	// Verify balances were not changed (rollback occurred)
	currentSenderBalance, err := getBalance(ctx, db, senderID)
	require.NoError(t, err)
	assert.Equal(t, initialSenderBalance, currentSenderBalance)

	currentReceiverBalance, err := getBalance(ctx, db, receiverID)
	require.NoError(t, err)
	assert.Equal(t, initialReceiverBalance, currentReceiverBalance)
}

type TransactionAuthorizerGatewayMock struct {
	authorize bool
}

func NewMockTransactionAuthorizerGateway(authorize bool) *TransactionAuthorizerGatewayMock {
	return &TransactionAuthorizerGatewayMock{
		authorize: authorize,
	}
}

func (m *TransactionAuthorizerGatewayMock) IsTransactionAllowed(ctx context.Context) bool {
	return m.authorize
}
