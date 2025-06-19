package handler

import (
	"context"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"log"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
)

type handler struct {
	createTransaction ICreateTransaction
	createUser        ICreateUser
	otel              telemetry.Telemetry
	logger            *log.Logger
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
	telemetry telemetry.Telemetry,
) *handler {
	return &handler{
		createTransaction: createTransaction,
		createUser:        createUser,
		otel:              telemetry,
		logger:            log.New(log.Writer(), "handler: ", log.LstdFlags),
	}
}
