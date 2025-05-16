package usecase_test

import (
	"context"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/app/usecase/strategy"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Execute_ShouldReturnErrorWhenEmailAlreadyExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockUserRepository{}
	mockStrategy := &mockCreateUserStrategy{}

	mockRepo.On("ExistsByEmail", ctx, mock.AnythingOfType("string")).Return(true, nil)

	createUser := usecase.NewCreateUser(mockRepo, []usecase.CreateUserStrategy{mockStrategy})

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "12345678901",
		UserType: "common",
	}

	// Act
	result, err := createUser.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, errs.ErrEmailAlreadyRegistered)
	mockRepo.AssertCalled(t, "ExistsByEmail", ctx, input.Email)
	mockStrategy.AssertNotCalled(t, "Execute")
}

func TestCreateUser_Execute_ShouldReturnErrUserTypeNotFoundWhenNoMatchingStrategyIsFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockUserRepository{}
	mockStrategy := &mockCreateUserStrategy{}

	mockRepo.On("ExistsByEmail", ctx, mock.AnythingOfType("string")).Return(false, nil)
	mockStrategy.On("UserType").Return("invalid_type")

	createUser := usecase.NewCreateUser(mockRepo, []usecase.CreateUserStrategy{mockStrategy})

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "12345678901",
		UserType: "merchant", // This type doesn't match any strategy
	}

	// Act
	result, err := createUser.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, errs.ErrUserTypeNotFound)
	mockRepo.AssertCalled(t, "ExistsByEmail", ctx, input.Email)
	mockStrategy.AssertNotCalled(t, "Execute")
}

func TestCreateUser_Execute_ShouldExecuteCorrectStrategyWhenMatchingUserTypeIsFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockUserRepository{}
	mockCommonStrategy := &mockCreateUserStrategy{}
	mockMerchantStrategy := &mockCreateUserStrategy{}

	mockRepo.On("ExistsByEmail", ctx, mock.AnythingOfType("string")).Return(false, nil)
	mockCommonStrategy.On("UserType").Return("common")
	mockMerchantStrategy.On("UserType").Return("merchant")
	mockMerchantStrategy.On("Execute", ctx, mock.AnythingOfType("strategy.CreateUserStrategyInput")).Return("user-123", nil)

	createUser := usecase.NewCreateUser(mockRepo, []usecase.CreateUserStrategy{mockCommonStrategy, mockMerchantStrategy})

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "50379007000134",
		UserType: "merchant",
	}

	// Act
	result, err := createUser.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "user-123", result)
	mockRepo.AssertCalled(t, "ExistsByEmail", ctx, input.Email)
	mockCommonStrategy.AssertNotCalled(t, "Execute")
	mockMerchantStrategy.AssertCalled(t, "Execute", ctx, mock.MatchedBy(func(input strategy.CreateUserStrategyInput) bool {
		return input.Name == "John Doe" &&
			input.Email == "john@example.com" &&
			input.Password == "password123" &&
			input.Document == "50379007000134"
	}))
}

type mockCreateUserStrategy struct {
	mock.Mock
}

func (m *mockCreateUserStrategy) Execute(ctx context.Context, input strategy.CreateUserStrategyInput) (string, error) {
	m.Called(ctx, input)
	return m.Called(ctx, input).String(0), m.Called(ctx, input).Error(1)
}

func (m *mockCreateUserStrategy) UserType() string {
	return m.Called().Get(0).(string)
}

func (m *mockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	m.Mock.Called(ctx, email)
	return m.Called(ctx, email).Get(0).(bool), m.Called(ctx, email).Error(1)
}
