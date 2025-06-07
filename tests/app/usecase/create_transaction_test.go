package usecase_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	repository "github.com.br/gibranct/simplified-wallet/internal/provider/repo"
	test "github.com.br/gibranct/simplified-wallet/tests"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

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
	container, db, err := test.SetupTestDatabase(ctx, migrateVersion)
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
	snsService := NewMockQueue()

	// Create use case
	createTransactionUseCase := usecase.NewCreateTransaction(userRepo, authorizerGateway, snsService)

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
	container, db, err := test.SetupTestDatabase(ctx, usersTableVersion)
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
	snsService := NewMockQueue()

	// Create use case with the failing repository
	createTransactionUseCase := usecase.NewCreateTransaction(userRepo, authorizerGateway, snsService)

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

type QueueMock struct{}

func (m *QueueMock) Send(ctx context.Context, message []byte) error {
	return nil
}

func NewMockQueue() *QueueMock {
	return &QueueMock{}
}
