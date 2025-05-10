package entity_test

import (
	"testing"

	domain "github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction_ShouldCreateTransactionWithProvidedPositiveAmount(t *testing.T) {
	// Arrange
	amount := 100.50
	senderID := "sender123"
	receiverID := "receiver456"

	// Act
	transaction, err := domain.NewTransaction(amount, senderID, receiverID)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.NotEmpty(t, transaction.ID())
	assert.Equal(t, int64(amount*100), transaction.Amount()) // Converted to pennies
	assert.Equal(t, senderID, transaction.SenderID())
	assert.Equal(t, receiverID, transaction.ReceiverID())
}
