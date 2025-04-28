package vo_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func Test_ValidCNPJs(t *testing.T) {
	tests := []string{
		"13.347.016/0001-17",
		"49.635.893/0001-30",
		"70.070.418/0001-50",
		"13347016000117",
		"59812745000106",
		"63533402000171",
	}

	for _, test := range tests {
		cnpj, err := vo.NewCNPJ(test)
		assert.Nil(t, err)
		assert.NotEmpty(t, cnpj.GetValue())
	}
}

func Test_InvalidCNPJs(t *testing.T) {
	tests := []string{
		"13.347.016/0001-18", // Invalid check digit
		"92.122.648/0001-00", // Invalid check digit
		"63.137.118/0001",    // Too short
		"63137118000196000",  // Too long
		"00000000000000",     // All digits the same
		"11111111111111",     // All digits the same
		"ABCDEFGHIJKLMN",     // Non-numeric
		"",                   // Empty string
	}

	for _, test := range tests {
		cnpj, err := vo.NewCNPJ(test)
		assert.NotNil(t, err)
		assert.Nil(t, cnpj)
		assert.Equal(t, vo.ErrInvalidCNPJ, err)
	}
}

func Test_CNPJGetValue(t *testing.T) {
	// Arrange
	validCNPJ := "13347016000117"

	// Act
	cnpj, err := vo.NewCNPJ(validCNPJ)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, validCNPJ, cnpj.GetValue())
}

func Test_CNPJWithSpecialCharacters(t *testing.T) {
	// Arrange
	formattedCNPJ := "13.347.016/0001-17"
	plainCNPJ := "13347016000117"

	// Act
	cnpj, err := vo.NewCNPJ(formattedCNPJ)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, plainCNPJ, cnpj.GetValue())
}
