package handler

import (
	"net/http"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	Amount     float64 `json:"amount"`
	SenderID   string  `json:"sender_id"`
	ReceiverID string  `json:"receiver_id"`
}

func (h handler) PostTransaction(w http.ResponseWriter, r *http.Request) {
	var input CreateTransactionRequest

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.writeJson(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	senderID, err := uuid.Parse(input.SenderID)
	if err != nil {
		h.writeJson(w, http.StatusBadRequest, envelope{"error": "invalid sender_id"}, nil)
		return
	}

	receiverID, err := uuid.Parse(input.ReceiverID)
	if err != nil {
		h.writeJson(w, http.StatusBadRequest, envelope{"error": "invalid receiver_id"}, nil)
		return
	}

	transactionID, err := h.createTransaction.Execute(r.Context(), usecase.CreateTransactionInput{
		Amount:     input.Amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	})

	if err != nil {
		h.writeJson(w, http.StatusUnprocessableEntity, envelope{"error": err.Error()}, nil)
		return
	}

	h.writeJson(w, http.StatusCreated, envelope{"transaction_id": transactionID}, nil)

	return
}
