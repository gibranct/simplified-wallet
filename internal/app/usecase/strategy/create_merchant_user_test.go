package strategy_test

import (
	"context"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase/strategy"
	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMerchantUser_Execute_ShouldReturnErrorWhenCNPJAlreadyExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockCreateMerchantUserRepository{}

	mockRepo.On("ExistsByCNPJ", ctx, mock.AnythingOfType("string")).Return(true, nil)

	createMerchantUser := strategy.NewCreateMerchantUser(mockRepo, telemetry.NewMockTelemetry())

	input := strategy.CreateUserStrategyInput{
		Name:     "Merchant Inc.",
		Email:    "merchant@example.com",
		Password: "securepass123",
		Document: "88529579000125",
	}

	// Act
	result, err := createMerchantUser.Execute(ctx, input)

	// Assert
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, errs.ErrCNPJAlreadyRegistered)
	mockRepo.AssertCalled(t, "ExistsByCNPJ", ctx, input.Document)
	mockRepo.AssertNotCalled(t, "Save")
}

func TestCreateMerchantUser_Execute_ShouldSuccessfullyCreateMerchantUserWhenAllInputIsValidAndUserDoesNotExist(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockCreateMerchantUserRepository{}

	mockRepo.On("ExistsByCNPJ", ctx, mock.AnythingOfType("string")).Return(false, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	createMerchantUser := strategy.NewCreateMerchantUser(mockRepo, telemetry.NewMockTelemetry())

	input := strategy.CreateUserStrategyInput{
		Name:     "Acme Corp",
		Email:    "acme@example.com",
		Password: "securepass123",
		Document: "88529579000125",
	}

	// Act
	result, err := createMerchantUser.Execute(ctx, input)

	// Assert
	assert.NotEmpty(t, result)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "ExistsByCNPJ", ctx, input.Document)
	mockRepo.AssertCalled(t, "Save", ctx, mock.AnythingOfType("*entity.User"))

	// Verify the created user
	mockRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user *entity.User) bool {
		return user.Name() == input.Name &&
			user.Email() == input.Email &&
			user.Type() == vo.MerchantUserType &&
			user.CNPJ() == input.Document
	}))
}

type mockCreateMerchantUserRepository struct {
	mock.Mock
}

func (m *mockCreateMerchantUserRepository) ExistsByCNPJ(ctx context.Context, cpf string) (bool, error) {
	args := m.Called(ctx, cpf)
	return args.Bool(0), args.Error(1)
}

func (m *mockCreateMerchantUserRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
