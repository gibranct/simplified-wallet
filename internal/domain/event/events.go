package event

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type CreateTransactionEventV1 struct {
	PublishedAt   string
	TransactionID string
	Amount        float64
	SenderID      uuid.UUID
	ReceiverID    uuid.UUID
}

func NewCreateTransactionEventV1(transactionID string, amount float64, senderID uuid.UUID, receiverID uuid.UUID) *CreateTransactionEventV1 {
	publishedAt := time.Now().Format(time.RFC3339)
	return &CreateTransactionEventV1{
		PublishedAt:   publishedAt,
		TransactionID: transactionID,
		Amount:        amount,
		SenderID:      senderID,
		ReceiverID:    receiverID,
	}
}

func (e *CreateTransactionEventV1) ToJSON() []byte {
	jsonData, err := json.Marshal(e)
	if err != nil {
		log.Printf("Error marshalling event to JSON: %v", err)
		return nil
	}
	return jsonData
}
