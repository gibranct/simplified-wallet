package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction_Execute_ShouldReturnErrorWhenTransactionIsNotAllowedByAuthorizer(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 100.0

	// Setup mock to deny transaction authorization
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(false)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, errs.ErrTransactionNotAllowed)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertNotCalled(t, "UpdateBalance")
}

func TestCreateTransaction_Execute_ShouldSuccessfullyUpdateBalanceAndCreateTransactionWhenAllConditionsAreValid(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 100.0
	expectedTransactionID := "transaction-id"

	// Mock users
	sender := NewUser(vo.CommonUserType)
	err := sender.Deposit(amount * 2)
	assert.NoError(t, err)
	receiver := NewUser(vo.CommonUserType)

	// Configure mocks
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)")).
		Run(func(args mock.Arguments) {
			// Execute the callback function to simulate the balance update
			updateFn := args.Get(3).(func(*entity.User, *entity.User) (*entity.Transaction, error))
			transaction, _ := updateFn(sender, receiver)
			// Set the transaction ID for verification
			if transaction != nil {
				expectedTransactionID = transaction.ID()
			}
		}).
		Return(nil)
	mockQueue.On("Send", ctx, mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, expectedTransactionID, result)
	assert.NoError(t, err)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)"))
	mockQueue.AssertCalled(t, "Send", ctx, mock.AnythingOfType("[]uint8"))
}

func TestCreateTransaction_Execute_ShouldReturnErrorWhenSenderIsAMerchant(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 100.0

	// Mock users with insufficient balance
	sender := NewUser(vo.MerchantUserType)
	// Note: Not depositing enough money (only 50)
	err := sender.Deposit(50.0)
	assert.NoError(t, err)
	receiver := NewUser(vo.CommonUserType)

	// Configure mocks
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)")).
		Run(func(args mock.Arguments) {
			// Execute the callback function to simulate the balance update
			updateFn := args.Get(3).(func(*entity.User, *entity.User) (*entity.Transaction, error))
			// This will fail because merchants cannot send money
			_, err := updateFn(sender, receiver)
			assert.ErrorIs(t, err, errs.ErrMerchantCannotSendMoney)
		}).
		Return(errs.ErrMerchantCannotSendMoney)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, errs.ErrMerchantCannotSendMoney)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)"))
}

func TestCreateTransaction_Execute_ShouldReturnErrorWhenSenderHasInsufficientFunds(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 100.0

	// Mock users with insufficient balance
	sender := NewUser(vo.CommonUserType)
	// Note: Not depositing enough money (only 50)
	err := sender.Deposit(50.0)
	assert.NoError(t, err)
	receiver := NewUser(vo.CommonUserType)

	// Configure mocks
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)")).
		Run(func(args mock.Arguments) {
			// Execute the callback function to simulate the balance update
			updateFn := args.Get(3).(func(*entity.User, *entity.User) (*entity.Transaction, error))
			// This will fail due to insufficient funds
			_, _ = updateFn(sender, receiver)
		}).
		Return(errs.ErrNotEnoughMoney)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, errs.ErrNotEnoughMoney)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)"))
}

func TestCreateTransaction_Execute_ShouldCreateTransactionWithCorrectAmountAndIDs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	amount := 150.0

	var capturedTransaction *entity.Transaction

	// Mock users
	sender := NewUser(vo.CommonUserType)
	err := sender.Deposit(amount)
	assert.NoError(t, err)
	receiver := NewUser(vo.CommonUserType)

	senderID := uuid.MustParse(sender.ID())
	receiverID := uuid.MustParse(receiver.ID())

	// Configure mocks
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)")).
		Run(func(args mock.Arguments) {
			// Execute the callback function and capture the transaction
			updateFn := args.Get(3).(func(*entity.User, *entity.User) (*entity.Transaction, error))
			transaction, _ := updateFn(sender, receiver)
			capturedTransaction = transaction
		}).
		Return(nil)
	mockQueue.On("Send", ctx, mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.NotNil(t, capturedTransaction)
	assert.Equal(t, int64(amount*100), capturedTransaction.Amount())
	assert.Equal(t, senderID.String(), capturedTransaction.SenderID())
	assert.Equal(t, receiverID.String(), capturedTransaction.ReceiverID())
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)"))
}

func TestCreateTransaction_Execute_ShouldIncreaseReceiverBalanceByCorrectAmount(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	senderID := uuid.New()
	receiverID := uuid.New()
	initialSenderBalance := 200.0
	initialReceiverBalance := 50.0
	amount := 100.0
	expectedTransactionID := "transaction-id"

	// Create users with initial balances
	sender := NewUser(vo.CommonUserType)
	err := sender.Deposit(initialSenderBalance)
	assert.NoError(t, err)
	receiver := NewUser(vo.CommonUserType)
	err = receiver.Deposit(initialReceiverBalance)
	assert.NoError(t, err)

	// Track balance changes
	var actualSenderBalance, actualReceiverBalance int64

	// Configure mocks
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)")).
		Run(func(args mock.Arguments) {
			// Execute the callback function to simulate the balance update
			updateFn := args.Get(3).(func(*entity.User, *entity.User) (*entity.Transaction, error))
			transaction, _ := updateFn(sender, receiver)
			// Capture the balances after transaction
			actualSenderBalance = sender.Balance()
			actualReceiverBalance = receiver.Balance()
			// Set the transaction ID for verification
			if transaction != nil {
				expectedTransactionID = transaction.ID()
			}
		}).
		Return(nil)
	mockQueue.On("Send", ctx, mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, expectedTransactionID, result)
	assert.NoError(t, err)
	assert.Equal(t, int64((initialSenderBalance-amount)*100), actualSenderBalance, "Sender balance should be decreased by transfer amount")
	assert.Equal(t, int64((initialReceiverBalance+amount)*100), actualReceiverBalance, "Receiver balance should be increased by exact transfer amount")
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)"))
}

func TestCreateTransaction_Execute_ShouldVerifySenderBalanceIsDecreasedByCorrectAmount(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	senderID := uuid.New()
	receiverID := uuid.New()
	initialBalance := 200.0
	amount := 75.0
	expectedRemainingBalance := int64((initialBalance - amount) * 100)

	// Create mock sender with initial balance
	sender := NewUser(vo.CommonUserType)
	err := sender.Deposit(initialBalance)
	assert.NoError(t, err)
	receiver := NewUser(vo.MerchantUserType)

	// Configure mocks
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)")).
		Run(func(args mock.Arguments) {
			// Execute the callback function to simulate the balance update
			updateFn := args.Get(3).(func(*entity.User, *entity.User) (*entity.Transaction, error))
			_, _ = updateFn(sender, receiver)
		}).
		Return(nil)
	mockQueue.On("Send", ctx, mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	_, err = useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRemainingBalance, sender.Balance())
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)"))
}

func TestCreateTransaction_Execute_ShouldHandleTransactionWithZeroAmount(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 0.0

	// Create mock users
	sender := NewUser(vo.CommonUserType)
	err := sender.Deposit(100.0) // Add some initial balance
	assert.NoError(t, err)
	receiver := NewUser(vo.CommonUserType)

	initialSenderBalance := sender.Balance()
	initialReceiverBalance := receiver.Balance()

	// Configure mocks
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)")).
		Run(func(args mock.Arguments) {
			// Execute the callback function to simulate the balance update
			updateFn := args.Get(3).(func(*entity.User, *entity.User) (*entity.Transaction, error))
			_, _ = updateFn(sender, receiver)
		}).
		Return(nil)
	mockQueue.On("Send", ctx, mock.AnythingOfType("[]uint8")).Return(nil)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.NotEmpty(t, result)
	assert.NoError(t, err)
	assert.Equal(t, initialSenderBalance, sender.Balance(), "Sender balance should not change")
	assert.Equal(t, initialReceiverBalance, receiver.Balance(), "Receiver balance should not change")
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)"))
	mockQueue.AssertCalled(t, "Send", ctx, mock.AnythingOfType("[]uint8"))
}

func TestCreateTransaction_Execute_ShouldPropagateErrorsFromRepositoryLayer(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := &mockUserRepository{}
	mockAuthorizer := &mockTransactionAuthorizerGateway{}
	mockQueue := &mockQueue{}

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 100.0
	expectedError := errors.New("database connection error")

	// Configure mocks
	mockAuthorizer.On("IsTransactionAllowed", ctx).Return(true)
	mockUserRepo.On("UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)")).
		Return(expectedError)

	useCase := usecase.NewCreateTransaction(mockUserRepo, mockAuthorizer, mockQueue)

	input := usecase.CreateTransactionInput{
		Amount:     amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, expectedError)
	mockAuthorizer.AssertCalled(t, "IsTransactionAllowed", ctx)
	mockUserRepo.AssertCalled(t, "UpdateBalance", ctx, senderID.String(), receiverID.String(), mock.AnythingOfType("func(*entity.User, *entity.User) (*entity.Transaction, error)"))
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) UpdateBalance(ctx context.Context, senderID, receiverID string, updateFn func(sender, receiver *entity.User) (*entity.Transaction, error)) error {
	args := m.Called(ctx, senderID, receiverID, updateFn)
	return args.Error(0)
}

type mockTransactionAuthorizerGateway struct {
	mock.Mock
}

func (m *mockTransactionAuthorizerGateway) IsTransactionAllowed(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func NewUser(userType string) *entity.User {
	var cpf, cnpj string

	if userType == vo.CommonUserType {
		cpf = "12345678901"
	} else if userType == vo.MerchantUserType {
		cnpj = "47775767000156"
	}

	user, err := entity.NewUser(
		faker.Name(),
		faker.Email(),
		faker.Password(),
		cpf,
		cnpj,
		userType,
	)
	if err != nil {
		panic(err)
	}
	return user
}

type mockQueue struct {
	mock.Mock
}

func (m *mockQueue) Send(ctx context.Context, message []byte) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}
