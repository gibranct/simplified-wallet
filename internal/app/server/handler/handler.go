package handler

import (
	"context"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
)

type handler struct {
	createTransaction ICreateTransaction
}

type ICreateTransaction interface {
	Execute(ctx context.Context, input usecase.CreateTransactionInput) (string, error)
}

func New(createTransaction ICreateTransaction) *handler {
	return &handler{
		createTransaction: createTransaction,
	}
}
