package vo_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestNewName_ShouldReturnErrorWhenLessThan3Characters(t *testing.T) {
	// Arrange
	tooShortName := "ab"

	// Act
	name, err := vo.NewName(tooShortName)

	// Assert
	assert.Nil(t, name)
	assert.ErrorIs(t, err, errs.ErrNameLength)
}

func TestNewName_ShouldReturnErrorWhenNameIsMoreThan50Characters(t *testing.T) {
	// Arrange
	longName := "ThisIsAReallyLongNameThatExceedsFiftyCharactersInLength123456789"

	// Act
	name, err := vo.NewName(longName)

	// Assert
	assert.Nil(t, name)
	assert.ErrorIs(t, err, errs.ErrNameLength)
}

func TestNewName_ShouldReturnValidNameObjectWhenNameIsExactly3Characters(t *testing.T) {
	// Arrange
	nameValue := "Abc"

	// Act
	name, err := vo.NewName(nameValue)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, name)
	assert.Equal(t, nameValue, name.Value())
}

func TestNewName_ShouldReturnValidNameObjectWhenNameIsExactly50Characters(t *testing.T) {
	// Arrange
	name := "ThisIsANameThatIsExactlyFiftyCharactersLongToTest1"

	// Act
	result, err := vo.NewName(name)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, name, result.Value())
}
