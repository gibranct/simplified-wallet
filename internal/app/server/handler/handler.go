package handler

import (
	"context"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
)

type handler struct {
	createTransaction ICreateTransaction
	createUser        ICreateUser
}

type ICreateTransaction interface {
	Execute(ctx context.Context, input usecase.CreateTransactionInput) (string, error)
}

type ICreateUser interface {
	Execute(ctx context.Context, input usecase.CreateUserInput) (string, error)
}

func New(
	createTransaction ICreateTransaction,
	createUser ICreateUser,
) *handler {
	return &handler{
		createTransaction: createTransaction,
		createUser:        createUser,
	}
}
