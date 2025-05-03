package vo_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestNewMoney_ShouldReturnErrorWhenAmountIsNegative(t *testing.T) {
	// Arrange
	negativeAmount := -1.0

	// Act
	money, err := vo.NewMoney(negativeAmount)

	// Assert
	assert.Nil(t, money)
	assert.ErrorIs(t, err, errs.ErrZeroOrNegativeAmount)
}

func TestNewMoney_ShouldCorrectlyConvertPositiveAmountsToPennies(t *testing.T) {
	// Arrange
	testCases := []struct {
		amount        float64
		expectedValue int64
	}{
		{amount: 0.0, expectedValue: 0},
		{amount: 10.50, expectedValue: 1050},
		{amount: 0.01, expectedValue: 1},
		{amount: 123.45, expectedValue: 12345},
		{amount: 999.99, expectedValue: 99999},
	}

	for _, tc := range testCases {
		// Act
		money, err := vo.NewMoney(tc.amount)

		// Assert
		assert.NotNil(t, money)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedValue, money.Value())
	}
}

func TestMoney_Subtract_ShouldReturnErrZeroOrNegativeAmountWhenSubtractingNegativeAmount(t *testing.T) {
	// Arrange
	initialAmount := 100.0
	money, _ := vo.NewMoney(initialAmount)
	negativeAmount := -20.0

	// Act
	result, err := money.Subtract(negativeAmount)

	// Assert
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errs.ErrZeroOrNegativeAmount)
	// Verify original money object remains unchanged
	assert.Equal(t, int64(10000), money.Value())
}

func TestMoney_Subtract_ShouldCorrectlySubtractSmallerAmountFromCurrentBalance(t *testing.T) {
	// Arrange
	initialAmount := 100.0
	money, _ := vo.NewMoney(initialAmount)
	amountToSubtract := 20.0
	expectedValue := int64(8000) // 100.0 - 20.0 = 80.0, converted to pennies = 8000

	// Act
	result, err := money.Subtract(amountToSubtract)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedValue, result.Value())
	// Verify original money object remains unchanged
	assert.Equal(t, int64(10000), money.Value())
}

func TestMoney_Add_ShouldCorrectlyAddPositiveAmountToCurrentBalance(t *testing.T) {
	// Arrange
	initialAmount := 100.0
	money, _ := vo.NewMoney(initialAmount)
	amountToAdd := 20.0
	expectedValue := int64(12000) // 100.0 + 20.0 = 120.0, converted to pennies = 12000

	// Act
	result, err := money.Add(amountToAdd)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedValue, result.Value())
	// Verify original money object remains unchanged
	assert.Equal(t, int64(10000), money.Value())
}
