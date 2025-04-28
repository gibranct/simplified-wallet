package vo_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)

func Test_ValidCPFs(t *testing.T) {
	tests := []string{"97456321558", "71428793860", "87748248800"}
	for _, test := range tests {
		cpf, err := vo.NewCPF(test)
		assert.Nil(t, err)
		assert.NotEmpty(t, cpf.GetValue())
	}
}

func Test_InvalidCPFs(t *testing.T) {
	tests := []string{"9745632155", "11111111111", "97a56321558"}
	for _, test := range tests {
		cpf, err := vo.NewCPF(test)
		assert.NotNil(t, err)
		assert.Nil(t, cpf)
	}
}
