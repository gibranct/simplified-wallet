package strategy_test

import (
	"context"
	"errors"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase/strategy"
	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCommonUser_Execute_ShouldReturnErrorWhenExistsByCPFFails(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockCreateCommonUserRepository{}
	expectedError := errors.New("database error")

	mockRepo.On("ExistsByCPF", ctx, mock.AnythingOfType("string")).Return(false, expectedError)

	createCommonUser := strategy.NewCreateCommonUser(mockRepo, telemetry.NewMockTelemetry())

	input := strategy.CreateUserStrategyInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "12345678901",
	}

	// Act
	result, err := createCommonUser.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, expectedError)
	mockRepo.AssertCalled(t, "ExistsByCPF", ctx, input.Document)
	mockRepo.AssertNotCalled(t, "Save")
}

func TestCreateCommonUser_Execute_ShouldReturnErrCPFAlreadyRegisteredWhenUserWithGivenCPFExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockCreateCommonUserRepository{}

	mockRepo.On("ExistsByCPF", ctx, mock.AnythingOfType("string")).Return(true, nil)

	createCommonUser := strategy.NewCreateCommonUser(mockRepo, telemetry.NewMockTelemetry())

	input := strategy.CreateUserStrategyInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "12345678901",
	}

	// Act
	result, err := createCommonUser.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, errs.ErrCPFAlreadyRegistered)
	mockRepo.AssertCalled(t, "ExistsByCPF", ctx, input.Document)
	mockRepo.AssertNotCalled(t, "Save")
}

func TestCreateCommonUser_Execute_ShouldReturnUserIDWhenUserIsSuccessfullyCreated(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockCreateCommonUserRepository{}

	mockRepo.On("ExistsByCPF", ctx, mock.AnythingOfType("string")).Return(false, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	createCommonUser := strategy.NewCreateCommonUser(mockRepo, telemetry.NewMockTelemetry())

	input := strategy.CreateUserStrategyInput{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "securepass123",
		Document: "98765432100",
	}

	// Act
	result, err := createCommonUser.Execute(ctx, input)

	// Assert
	assert.NotEmpty(t, result)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "ExistsByCPF", ctx, input.Document)
	mockRepo.AssertCalled(t, "Save", ctx, mock.AnythingOfType("*entity.User"))
}

func TestCreateCommonUser_Execute_ShouldCreateUserWithCommonUserType(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockCreateCommonUserRepository{}

	mockRepo.On("ExistsByCPF", ctx, mock.AnythingOfType("string")).Return(false, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*entity.User)
		assert.Equal(t, vo.CommonUserType, user.Type())
	})

	createCommonUser := strategy.NewCreateCommonUser(mockRepo, telemetry.NewMockTelemetry())

	input := strategy.CreateUserStrategyInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "12345678901",
	}

	// Act
	result, err := createCommonUser.Execute(ctx, input)

	// Assert
	assert.NotEmpty(t, result)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "ExistsByCPF", ctx, input.Document)
	mockRepo.AssertCalled(t, "Save", ctx, mock.AnythingOfType("*entity.User"))
}

type mockCreateCommonUserRepository struct {
	mock.Mock
}

func (m *mockCreateCommonUserRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockCreateCommonUserRepository) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	args := m.Called(ctx, cpf)
	return args.Get(0).(bool), args.Error(1)
}
