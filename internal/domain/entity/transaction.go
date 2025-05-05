package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	id         uuid.UUID
	amount     float64
	senderID   string
	receiverID string
	createdAt  time.Time
}

func (t *Transaction) ID() string {
	return t.id.String()
}

func (t *Transaction) Amount() float64 {
	return t.amount
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

func NewTransaction(amount float64, senderID, receiverID string) *Transaction {
	id := uuid.New()

	transaction := &Transaction{
		id:         id,
		amount:     amount,
		senderID:   senderID,
		receiverID: receiverID,
		createdAt:  time.Now(),
	}

	return transaction
}
