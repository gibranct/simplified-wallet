package usecase_test

import (
	"context"
	"database/sql"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"log"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/app/usecase/strategy"
	repository "github.com.br/gibranct/simplified-wallet/internal/provider/repo"
	test "github.com.br/gibranct/simplified-wallet/tests"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCommonUser_Integration_Success(t *testing.T) {
	ctx := context.Background()
	migrateVersion := uint(2) // Use the latest migration version

	// Setup
	container, db, err := test.SetupTestDatabase(ctx, migrateVersion)
	require.NoError(t, err)
	otel, err := telemetry.NewJaeger(context.Background(), "")
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer db.Close()
	defer otel.Shutdown(ctx)

	userRepo := repository.NewUserRepository(db, otel)

	// Create use case
	createUserUseCase := usecase.NewCreateUser(userRepo, strategies(userRepo), otel)

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "12345678901",
		UserType: "common",
	}

	// Act
	userID, err := createUserUseCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	// Verify common user was recorded
	var id, name, email, cpf, user_type string
	var cnpj sql.NullString
	err = db.QueryRowContext(ctx, "SELECT id,name,email,cpf,cnpj,user_type FROM users WHERE id = $1", userID).Scan(&id, &name, &email, &cpf, &cnpj, &user_type)
	require.NoError(t, err)
	v, err := cnpj.Value()
	require.NoError(t, err)
	assert.Equal(t, "John Doe", name)
	assert.Equal(t, "john@example.com", email)
	assert.Equal(t, "12345678901", cpf)
	assert.Nil(t, v)
	assert.Equal(t, "common", user_type)
}

func TestCreateMerchantUser_Integration_Success(t *testing.T) {
	ctx := context.Background()
	migrateVersion := uint(2) // Use the latest migration version

	// Setup
	container, db, err := test.SetupTestDatabase(ctx, migrateVersion)
	require.NoError(t, err)
	otel, err := telemetry.NewJaeger(context.Background(), "")
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer db.Close()
	defer otel.Shutdown(ctx)

	userRepo := repository.NewUserRepository(db, otel)

	// Create use case
	createUserUseCase := usecase.NewCreateUser(userRepo, strategies(userRepo), otel)

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john2@example.com",
		Password: "password123",
		Document: "13521579000180",
		UserType: "merchant",
	}

	// Act
	userID, err := createUserUseCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	// Verify common user was recorded
	var id, name, email, cnpj, user_type string
	var cpf sql.NullString
	err = db.QueryRowContext(ctx, "SELECT id,name,email,cnpj,cpf,user_type FROM users WHERE id = $1", userID).Scan(&id, &name, &email, &cnpj, &cpf, &user_type)
	require.NoError(t, err)
	v, err := cpf.Value()
	require.NoError(t, err)
	assert.Equal(t, "John Doe", name)
	assert.Equal(t, "john2@example.com", email)
	assert.Equal(t, "13521579000180", cnpj)
	assert.Equal(t, "merchant", user_type)
	assert.Nil(t, v)
}

func TestCreateUser_ShouldFailIfEmailIsAlreadyRegistered(t *testing.T) {
	ctx := context.Background()
	migrateVersion := uint(3) // Use the latest migration version

	// Setup
	container, db, err := test.SetupTestDatabase(ctx, migrateVersion)
	require.NoError(t, err)
	otel, err := telemetry.NewJaeger(context.Background(), "")
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer db.Close()
	defer otel.Shutdown(ctx)

	userRepo := repository.NewUserRepository(db, otel)

	// Create use case
	createUserUseCase := usecase.NewCreateUser(userRepo, strategies(userRepo), otel)

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john2@example.com",
		Password: "password123",
		Document: "31094680001",
		UserType: "common",
	}

	// Act
	createUserUseCase.Execute(ctx, input)
	userID, err := createUserUseCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "email already registered", err.Error())
	assert.Empty(t, userID)
}

func TestCreateUser_ShouldFailIfCPFIsAlreadyRegistered(t *testing.T) {
	ctx := context.Background()
	migrateVersion := uint(3) // Use the latest migration version

	// Setup
	container, db, err := test.SetupTestDatabase(ctx, migrateVersion)
	require.NoError(t, err)
	otel, err := telemetry.NewJaeger(context.Background(), "")
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer db.Close()
	defer otel.Shutdown(ctx)

	userRepo := repository.NewUserRepository(db, otel)

	// Create use case
	createUserUseCase := usecase.NewCreateUser(userRepo, strategies(userRepo), otel)

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "12345678901",
		UserType: "common",
	}

	// Act
	createUserUseCase.Execute(ctx, input)
	userID, err := createUserUseCase.Execute(ctx, usecase.CreateUserInput{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "password456",
		Document: input.Document,
		UserType: "common",
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "cpf already registered", err.Error())
	assert.Empty(t, userID)
}

func TestCreateUser_ShouldFailIfCNPJIsAlreadyRegistered(t *testing.T) {
	ctx := context.Background()
	migrateVersion := uint(3) // Use the latest migration version

	// Setup
	container, db, err := test.SetupTestDatabase(ctx, migrateVersion)
	require.NoError(t, err)
	otel, err := telemetry.NewJaeger(context.Background(), "")
	require.NoError(t, err)
	defer container.Terminate(ctx)
	defer db.Close()
	defer otel.Shutdown(ctx)

	userRepo := repository.NewUserRepository(db, otel)

	// Create use case
	createUserUseCase := usecase.NewCreateUser(userRepo, strategies(userRepo), otel)

	input := usecase.CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Document: "01986061000132",
		UserType: "merchant",
	}

	// Act
	createUserUseCase.Execute(ctx, input)
	userID, err := createUserUseCase.Execute(ctx, usecase.CreateUserInput{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "password456",
		Document: input.Document,
		UserType: "merchant",
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "cnpj already registered", err.Error())
	assert.Empty(t, userID)
}

func strategies(userRepo repository.UserRepository) []usecase.CreateUserStrategy {
	otel, err := telemetry.NewJaeger(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}
	return []usecase.CreateUserStrategy{
		strategy.NewCreateCommonUser(userRepo, otel),
		strategy.NewCreateMerchantUser(userRepo, otel),
	}
}
