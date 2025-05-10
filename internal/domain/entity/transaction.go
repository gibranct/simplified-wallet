package entity

import (
	"time"

	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/google/uuid"
)

type Transaction struct {
	id         uuid.UUID
	amount     *vo.Money
	senderID   string
	receiverID string
	createdAt  time.Time
}

func (t *Transaction) ID() string {
	return t.id.String()
}

func (t *Transaction) Amount() int64 {
	return t.amount.Value()
}

func (t *Transaction) SenderID() string {
	return t.senderID
}

func (t *Transaction) ReceiverID() string {
	return t.receiverID
}

func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

func NewTransaction(amount float64, senderID, receiverID string) (*Transaction, error) {
	id := uuid.New()

	money, err := vo.NewMoney(amount)
	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		id:         id,
		amount:     money,
		senderID:   senderID,
		receiverID: receiverID,
		createdAt:  time.Now(),
	}

	return transaction, nil
}
