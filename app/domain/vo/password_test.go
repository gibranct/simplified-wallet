package vo_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/app/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestNewPassword_ShouldReturnErrPasswordTooShortWhenPasswordIsLessThan6Characters(t *testing.T) {
	// Arrange
	shortPassword := "12345"

	// Act
	password, err := vo.NewPassword(shortPassword)

	// Assert
	assert.Nil(t, password)
	assert.ErrorIs(t, err, vo.ErrPasswordTooShort)
}

func TestNewPassword_ShouldSuccessfullyCreatePasswordObjectWithValidPassword(t *testing.T) {
	// Arrange
	validPassword := "validPassword"

	// Act
	password, err := vo.NewPassword(validPassword)

	// Assert
	assert.NotNil(t, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, password.Value)
	// Check that the stored value is not the plain text password
	assert.NotEqual(t, validPassword, password.Value)
	// Verify password was hashed correctly by comparing
	assert.True(t, password.Compare(password.Value, validPassword))
}

func TestNewPassword_ShouldProperlyHashThePasswordUsingBcrypt(t *testing.T) {
	// Arrange
	plainPassword := "mySecurePassword123"

	// Act
	password, err := vo.NewPassword(plainPassword)

	// Assert
	assert.NotNil(t, password)
	assert.NoError(t, err)

	// Verify password was properly hashed
	// 1. Not equal to original
	assert.NotEqual(t, plainPassword, password.Value)

	// 2. Should start with bcrypt identifier $2a$ or similar
	assert.Contains(t, password.Value, "$2a$")

	// 3. Original password should compare correctly with hash
	assert.True(t, password.Compare(password.Value, plainPassword))

	// 4. Wrong password should not match the hash
	assert.False(t, password.Compare(password.Value, "wrongPassword"))
}
