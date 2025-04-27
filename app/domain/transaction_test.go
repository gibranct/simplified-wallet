package domain_test

import (
	"testing"

	"github.com.br/gibranct/simplified-wallet/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction_ShouldCreateTransactionWithProvidedPositiveAmount(t *testing.T) {
	// Arrange
	amount := 100.50
	senderID := "sender123"
	receiverID := "receiver456"

	// Act
	transaction := domain.NewTransaction(amount, senderID, receiverID)

	// Assert
	assert.NotEmpty(t, transaction.ID())
	assert.Equal(t, amount, transaction.Amount())
	assert.Equal(t, senderID, transaction.SenderID())
	assert.Equal(t, receiverID, transaction.ReceiverID())
}
