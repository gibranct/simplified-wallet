package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com.br/gibranct/simplified-wallet/app/domain"
	"github.com.br/gibranct/simplified-wallet/app/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction_Execute_ShouldReturnErrSenderNotFoundWhenTheSenderDoesNotExist(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockTransactionRepo := &mockTransactionRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}

	senderID := uuid.New()
	receiverID := uuid.New()

	// Setup the mock to return error for sender
	mockUserRepo.On("GetUserByID", ctx, senderID).Return(domain.User{}, errors.New("user not found"))

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockTransactionRepo, mockAuthorizer)

	input := usecase.CreateTransactionInput{
		Amount:     100.0,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, domain.ErrSenderNotFound)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, senderID)
	mockUserRepo.AssertNotCalled(t, "GetUserByID", ctx, receiverID)
	mockUserRepo.AssertNotCalled(t, "UpdateBalance")
	mockTransactionRepo.AssertNotCalled(t, "CreateTransaction")
	mockAuthorizer.AssertNotCalled(t, "IsTransactionAllowed")
}

func TestCreateTransaction_Execute_ShouldReturnErrReceiverNotFoundWhenTheReceiverDoesNotExist(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockTransactionRepo := &mockTransactionRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}

	senderID := uuid.New()
	receiverID := uuid.New()

	// Create a mock sender user
	sender := domain.User{}

	// Setup the mock to return sender successfully but error for receiver
	mockUserRepo.On("GetUserByID", ctx, senderID).Return(sender, nil)
	mockUserRepo.On("GetUserByID", ctx, receiverID).Return(domain.User{}, errors.New("user not found"))

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockTransactionRepo, mockAuthorizer)

	input := usecase.CreateTransactionInput{
		Amount:     100.0,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, domain.ErrReceiverNotFound)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, senderID)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, receiverID)
	mockUserRepo.AssertNotCalled(t, "UpdateBalance")
	mockTransactionRepo.AssertNotCalled(t, "CreateTransaction")
	mockAuthorizer.AssertNotCalled(t, "IsTransactionAllowed")
}

func TestCreateTransaction_Execute_ShouldReturnErrTransactionNotAllowedWhenTheTransactionIsNotAuthorized(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockTransactionRepo := &mockTransactionRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}

	senderID := uuid.New()
	receiverID := uuid.New()

	// Create mock users
	sender := domain.User{}
	receiver := domain.User{}

	// Setup the mocks to return users successfully but deny transaction authorization
	mockUserRepo.On("GetUserByID", ctx, senderID).Return(sender, nil)
	mockUserRepo.On("GetUserByID", ctx, receiverID).Return(receiver, nil)
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(false)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockTransactionRepo, mockAuthorizer)

	input := usecase.CreateTransactionInput{
		Amount:     100.0,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, domain.ErrTransactionNotAllowed)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, senderID)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, receiverID)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertNotCalled(t, "UpdateBalance")
	mockTransactionRepo.AssertNotCalled(t, "CreateTransaction")
}

func TestCreateTransaction_Execute_ShouldReturnErrNotEnoughMoneyWhenSenderHasInsufficientFunds(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockTransactionRepo := &mockTransactionRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}

	senderID := uuid.New()
	receiverID := uuid.New()

	// Create mock sender with insufficient balance
	sender := domain.User{}
	receiver := domain.User{}

	// Setup mocks
	mockUserRepo.On("GetUserByID", ctx, senderID).Return(sender, nil)
	mockUserRepo.On("GetUserByID", ctx, receiverID).Return(receiver, nil)
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, mock.Anything, mock.Anything, mock.AnythingOfType("func(*domain.User, *domain.User) (string, error)")).
		Return(domain.ErrNotEnoughMoney)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockTransactionRepo, mockAuthorizer)

	input := usecase.CreateTransactionInput{
		Amount:     100.0,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, domain.ErrNotEnoughMoney)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, senderID)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, receiverID)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, mock.Anything, mock.Anything, mock.AnythingOfType("func(*domain.User, *domain.User) (string, error)"))
	mockTransactionRepo.AssertNotCalled(t, "CreateTransaction")
}

func TestCreateTransaction_Execute_ShouldReturnErrorWhenTransactionRepositoryFailsToCreateTransaction(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockTransactionRepo := &mockTransactionRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}

	senderID := uuid.New()
	receiverID := uuid.New()

	// Create mock users
	sender := domain.User{}
	receiver := domain.User{}

	// Setup mocks
	mockUserRepo.On("GetUserByID", ctx, senderID).Return(sender, nil)
	mockUserRepo.On("GetUserByID", ctx, receiverID).Return(receiver, nil)
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)

	// Mock the UpdateBalance function to execute the callback with error handling
	mockUserRepo.On("UpdateBalance", ctx, mock.Anything, mock.Anything, mock.AnythingOfType("func(*domain.User, *domain.User) (string, error)")).
		Run(func(args mock.Arguments) {
			// Extract and execute the updateFn to trigger repo error
			updateFn := args.Get(3).(func(*domain.User, *domain.User) (string, error))

			// Provide dummy user objects to the callback
			senderObj := &domain.User{}
			senderObj.Deposit(100)
			receiverObj := &domain.User{}

			// The updateFn will call CreateTransaction internally
			updateFn(senderObj, receiverObj)
		}).
		Return(errors.New("transaction creation failed"))

	// Setup transaction repo mock to return error
	mockTransactionRepo.On("CreateTransaction", ctx, mock.AnythingOfType("domain.Transaction")).
		Return(domain.Transaction{}, errors.New("database error"))

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockTransactionRepo, mockAuthorizer)

	input := usecase.CreateTransactionInput{
		Amount:     100.0,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.Error(t, err)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, senderID)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, receiverID)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, mock.Anything, mock.Anything, mock.AnythingOfType("func(*domain.User, *domain.User) (string, error)"))
	mockTransactionRepo.AssertCalled(t, "CreateTransaction", ctx, mock.AnythingOfType("domain.Transaction"))
}

func TestCreateTransaction_Execute_ShouldSuccessfullyCreateATransactionWhenAllConditionsAreMet(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockTransactionRepo := &mockTransactionRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 100.0

	// Create mock users with methods implementation
	sender := domain.User{}
	sender.Deposit(amount)
	receiver := domain.User{}
	transaction := domain.NewTransaction(amount, senderID.String(), receiverID.String())
	transactionID := transaction.ID()

	// Setup mocks
	mockUserRepo.On("GetUserByID", ctx, senderID).Return(sender, nil)
	mockUserRepo.On("GetUserByID", ctx, receiverID).Return(receiver, nil)
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, mock.Anything, mock.Anything, mock.AnythingOfType("func(*domain.User, *domain.User) (string, error)")).
		Run(func(args mock.Arguments) {
			// Execute the callback function to simulate the update
			updateFn := args.Get(3).(func(*domain.User, *domain.User) (string, error))
			updateFn(&sender, &receiver)
		}).
		Return(nil)
	mockTransactionRepo.On("CreateTransaction", ctx, mock.AnythingOfType("domain.Transaction")).
		Return(transaction, nil)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockTransactionRepo, mockAuthorizer)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, transactionID, result)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), sender.Balance())
	assert.Equal(t, amount, receiver.Balance())
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, senderID)
	mockUserRepo.AssertCalled(t, "GetUserByID", ctx, receiverID)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, mock.Anything, mock.Anything, mock.AnythingOfType("func(*domain.User, *domain.User) (string, error)"))
	mockTransactionRepo.AssertCalled(t, "CreateTransaction", ctx, mock.AnythingOfType("domain.Transaction"))
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepository) UpdateBalance(ctx context.Context, senderID, receiverID string, updateFn func(sender, receiver *domain.User) (string, error)) error {
	args := m.Called(ctx, senderID, receiverID, updateFn)
	return args.Error(0)
}

type mockTransactionRepository struct {
	mock.Mock
}

func (m *mockTransactionRepository) CreateTransaction(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(domain.Transaction), args.Error(1)
}

type mockTransactionAuthorizerGateway struct {
	mock.Mock
}

func (m *mockTransactionAuthorizerGateway) IsTransactionAllowed(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}
