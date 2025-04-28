package vo_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func Test_CreateValidEmail(t *testing.T) {
	validEmails := []string{"john@doe.com", "gil@bil.com"}
	for _, n := range validEmails {
		email, err := vo.NewEmail(n)
		assert.Nil(t, err)
		assert.Equal(t, n, email.GetValue())
	}
}

func Test_CreateInvalidEmails(t *testing.T) {
	invalidEmails := []string{"johndoecom", "gilbil.com"}
	for _, n := range invalidEmails {
		email, err := vo.NewEmail(n)
		assert.NotNil(t, err)
		assert.Nil(t, email)
	}
}
