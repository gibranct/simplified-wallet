package vo_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/app/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestNewUserType_ReturnsValidUserType(t *testing.T) {
	// Act
	testCases := []struct {
		input    string
		expected string
	}{
		{"common", "common"},
		{"merchant", "merchant"},
	}

	// Assert
	for _, tc := range testCases {
		userType, err := vo.NewUserType(tc.input)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, userType.Value())
	}
}

func TestNewUserType_WithEmptyString_ReturnsError(t *testing.T) {
	// Act
	userType, err := vo.NewUserType("")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid user type: ", err.Error())
	assert.Equal(t, userType.Value(), "")
}
