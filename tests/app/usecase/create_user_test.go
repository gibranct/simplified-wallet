package test

import (
	"context"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	repository "github.com.br/gibranct/simplified-wallet/internal/provider/repo"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Integration_Success(t *testing.T) {
	ctx := context.Background()
	migrateVersion := uint(2) // Use the latest migration version

	// Setup
	container, db, err := setupTestDatabase(ctx, migrateVersion)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	// Create use case
	createUserUseCase := usecase.NewCreateUser(userRepo)

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		CPF:      "12345678901",
	}

	// Act
	userID, err := createUserUseCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	// Verify common user was recorded
	var id, name, email, cpf, user_type string
	err = db.QueryRowContext(ctx, "SELECT id,name,email,cpf,user_type FROM users WHERE id = $1", userID).Scan(&id, &name, &email, &cpf, &user_type)
	require.NoError(t, err)
	assert.Equal(t, "John Doe", name)
	assert.Equal(t, "john@example.com", email)
	assert.Equal(t, "12345678901", cpf)
	assert.Equal(t, "common", user_type)
}
