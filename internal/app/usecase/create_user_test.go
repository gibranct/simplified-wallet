package usecase_test

import (
	"context"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Execute_ShouldCreateUserWithValidInputAndReturnUserID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockUserRepository{}

	expectedUserID := "user-123"
	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		CPF:      "12345678901",
	}

	// Setup mock to capture the user being saved and return its ID
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*entity.User)
			// Verify user was created with correct data
			assert.Equal(t, input.Name, user.Name())
			assert.Equal(t, input.Email, user.Email())
			assert.Equal(t, input.CPF, user.CPF())
			assert.Equal(t, vo.CommonUserType, user.UserType())
			assert.Empty(t, user.CNPJ())

			// Simulate the user ID for verification
			expectedUserID = user.ID()
		}).
		Return(nil)

	useCase := usecase.NewCreateUser(mockRepo)

	// Act
	userID, err := useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
	mockRepo.AssertCalled(t, "Save", ctx, mock.AnythingOfType("*entity.User"))
}

func TestCreateUser_Execute_ShouldReturnErrorWhenUserCreationFailsDueToInvalidName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockUserRepository{}

	input := usecase.CreateUserInput{
		Name:     "", // Empty name to trigger validation error
		Email:    "john@example.com",
		Password: "password123",
		CPF:      "12345678901",
	}

	useCase := usecase.NewCreateUser(mockRepo)

	// Act
	userID, err := useCase.Execute(ctx, input)

	// Assert
	assert.Empty(t, userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name")
	mockRepo.AssertNotCalled(t, "Save") // Ensure Save was not called
}

func (m *mockUserRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil // Assume email uniqueness is not enforced
}

func (m *mockUserRepository) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	return false, nil // Assume CPF uniqueness is not enforced
}
