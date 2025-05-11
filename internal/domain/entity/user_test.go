package entity_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com/stretchr/testify/assert"
)

func TestNewUser_ShouldSuccessfullyCreateCommonUserWithValidInputParameters(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := "12345678909"
	cnpj := ""
	userType := "common"

	// Act
	user, err := entity.NewUser(name, email, password, cpf, cnpj, userType)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID())
	assert.Equal(t, name, user.Name())
	assert.Equal(t, email, user.Email())
	assert.NotEqual(t, password, user.Password()) // Password should be hashed
	assert.Equal(t, int64(0), user.Balance())
	assert.Equal(t, cpf, user.CPF())
	assert.Empty(t, user.CNPJ())
	assert.Equal(t, userType, user.UserType())
	assert.True(t, user.Active())
	assert.NotZero(t, user.CreatedAt())
	assert.NotZero(t, user.UpdatedAt())
}

func TestNewUser_ShouldSuccessfullyCreateMerchantUserWithValidInputParameters(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := ""
	cnpj := "85043353000121"
	userType := "merchant"

	// Act
	user, err := entity.NewUser(name, email, password, cpf, cnpj, userType)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID())
	assert.Equal(t, name, user.Name())
	assert.Equal(t, email, user.Email())
	assert.NotEqual(t, password, user.Password()) // Password should be hashed
	assert.Equal(t, int64(0), user.Balance())
	assert.Empty(t, user.CPF())
	assert.Equal(t, cnpj, user.CNPJ())
	assert.Equal(t, userType, user.UserType())
	assert.True(t, user.Active())
	assert.NotZero(t, user.CreatedAt())
	assert.NotZero(t, user.UpdatedAt())
}

func TestNewUser_ShouldReturnErrorWhenInvalidUserTypeIsProvided(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := "12345678909"
	cnpj := "12345678000190"
	invalidUserType := "invalid_type" // An invalid user type

	// Act
	_, err := entity.NewUser(name, email, password, cpf, cnpj, invalidUserType)

	// Assert
	assert.Error(t, err)
}

func TestNewUser_ShouldReturnErrorWhenWrongOrTooMuchDocumentIsProvided(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	testCases := []struct {
		expectedError error
		cpf           string
		cnpj          string
		userType      string
	}{
		{
			expectedError: errs.ErrCPFMustBeProvidedForCommonUser,
			cpf:           "",
			cnpj:          "12345678000190",
			userType:      "common",
		},
		{
			expectedError: errs.ErrCNPJMustBeProvidedForMerchant,
			cpf:           "12345678909",
			cnpj:          "",
			userType:      "merchant",
		},
		{
			expectedError: errs.ErrMerchantCannotHaveCPF,
			cpf:           "12345678909",
			cnpj:          "12345678000190",
			userType:      "merchant",
		},
		{
			expectedError: errs.ErrCommonCannotHaveCNPJ,
			cpf:           "12345678909",
			cnpj:          "12345678000190",
			userType:      "common",
		},
	}

	// Act & Assert
	for _, testCase := range testCases {
		_, err := entity.NewUser(name, email, password, testCase.cpf, testCase.cnpj, testCase.userType)
		assert.Error(t, err)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestUser_Deposit_ShouldCorrectlyUpdateBalanceWithPositiveAmount(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := "12345678909"
	cnpj := ""
	userType := "common"
	initialBalance := 0.0
	depositAmount := 100.0

	user, err := entity.NewUser(name, email, password, cpf, cnpj, userType)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), user.Balance())

	// Act
	user.Deposit(depositAmount)

	// Assert
	assert.Equal(t, int64(initialBalance+depositAmount)*100, user.Balance())
}

func TestUser_Deposit_ShouldNotChangeBalanceWithZeroAmount(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := "12345678909"
	cnpj := ""
	userType := "common"
	depositAmount := 0.0

	user, err := entity.NewUser(name, email, password, cpf, cnpj, userType)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), user.Balance())

	// Act
	user.Deposit(depositAmount)

	// Assert
	assert.Equal(t, int64(0), user.Balance())
}

func TestUser_Deposit_ShouldNotChangeBalanceWithNegativeAmount(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := "12345678909"
	cnpj := ""
	userType := "common"
	initialBalance := 0.0
	depositAmount := -50.0

	user, err := entity.NewUser(name, email, password, cpf, cnpj, userType)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), user.Balance())

	// Act
	user.Deposit(depositAmount)

	// Assert
	assert.Equal(t, int64(initialBalance*100), user.Balance())
}

func TestUser_Withdraw_ShouldSuccessfullyWithdrawWhenAmountIsPositiveAndLessThanBalance(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := "12345678909"
	cnpj := ""
	userType := "common"
	initialBalance := 100.0

	user, err := entity.NewUser(name, email, password, cpf, cnpj, userType)
	assert.NoError(t, err)

	// Set initial balance with deposit
	user.Deposit(initialBalance)
	assert.Equal(t, int64(initialBalance*100), user.Balance())

	withdrawAmount := 50.0

	// Act
	err = user.Withdraw(withdrawAmount)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(initialBalance-withdrawAmount)*100, user.Balance())
}

func TestUser_Withdraw_ShouldSuccessfullyWithdrawWhenAmountEqualsEntireBalance(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := "12345678909"
	cnpj := ""
	userType := "common"
	initialBalance := 100.0

	user, err := entity.NewUser(name, email, password, cpf, cnpj, userType)
	assert.NoError(t, err)

	// Set initial balance with deposit
	user.Deposit(initialBalance)
	assert.Equal(t, int64(initialBalance*100), user.Balance())

	withdrawAmount := 100.0 // Withdrawing the entire balance

	// Act
	err = user.Withdraw(withdrawAmount)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(0), user.Balance())
}

func TestUser_Withdraw_ShouldReturnFalseWhenTryingToWithdrawMoreThanAvailableBalance(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john@example.com"
	password := "validPassword123"
	cpf := "12345678909"
	cnpj := ""
	userType := "common"
	initialBalance := 50.0

	user, err := entity.NewUser(name, email, password, cpf, cnpj, userType)
	assert.NoError(t, err)

	// Set initial balance with deposit
	user.Deposit(initialBalance)
	assert.Equal(t, int64(initialBalance*100), user.Balance())

	withdrawAmount := 100.0 // Attempting to withdraw more than available

	// Act
	err = user.Withdraw(withdrawAmount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, int64(initialBalance*100), user.Balance()) // Balance should remain unchanged
}
