package handler

import (
	"github.com.br/gibranct/simplified-wallet/internal/provider/metrics"
	"go.opentelemetry.io/otel/attribute"
	"net/http"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com/google/uuid"
)

type PostTransactionRequest struct {
	Amount     float64 `json:"amount"`
	SenderID   string  `json:"sender_id"`
	ReceiverID string  `json:"receiver_id"`
}

func (h handler) PostTransaction(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.otel.Start(r.Context(), "PostTransaction")
	defer span.End()

	var input PostTransactionRequest

	err := h.readJSON(w, r, &input)
	if err != nil {
		err = h.writeJson(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		if err != nil {
			h.logger.Println(err)
		}
		return
	}

	senderID, err := uuid.Parse(input.SenderID)
	if err != nil {
		err = h.writeJson(w, http.StatusBadRequest, envelope{"error": "invalid sender_id"}, nil)
		if err != nil {
			h.logger.Println(err)
		}
		return
	}

	receiverID, err := uuid.Parse(input.ReceiverID)
	if err != nil {
		err = h.writeJson(w, http.StatusBadRequest, envelope{"error": "invalid receiver_id"}, nil)
		if err != nil {
			h.logger.Println(err)
		}
		return
	}

	transactionID, err := h.createTransaction.Execute(ctx, usecase.CreateTransactionInput{
		Amount:     input.Amount,
		SenderID:   senderID,
		ReceiverID: receiverID,
	})

	if err != nil {
		err = h.writeJson(w, http.StatusUnprocessableEntity, envelope{"error": err.Error()}, nil)
		if err != nil {
			h.logger.Println(err)
		}
		return
	}

	err = h.writeJson(w, http.StatusCreated, envelope{"transaction_id": transactionID}, nil)
	if err != nil {
		err = h.writeJson(w, http.StatusInternalServerError, envelope{"error": "failed to write response"}, nil)
		if err != nil {
			h.logger.Println(err)
		}
	}

	// After successful transaction creation:
	metrics.TransactionCounter.Inc()
	metrics.TransactionAmount.Add(input.Amount)

	// Add transaction details to the span
	span.SetAttributes(
		attribute.Float64("transaction.amount", input.Amount),
		attribute.String("transaction.sender_id", input.SenderID),
		attribute.String("transaction.receiver_id", input.ReceiverID),
	)
}
