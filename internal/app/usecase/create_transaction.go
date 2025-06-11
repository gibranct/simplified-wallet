package usecase

import (
	"context"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/event"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UpdateBalance(ctx context.Context, senderID, receiverID string, updateFn func(sender, receiver *entity.User) (*entity.Transaction, error)) error
}

type TransactionAuthorizerGateway interface {
	IsTransactionAllowed(ctx context.Context) bool
}

type Queue interface {
	Send(ctx context.Context, message []byte) error
}

type CreateTransaction struct {
	userRepository        UserRepository
	transactionAuthorizer TransactionAuthorizerGateway
	queue                 Queue
	otel                  telemetry.Telemetry
}
type CreateTransactionInput struct {
	Amount     float64
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
}

func (c *CreateTransaction) Execute(ctx context.Context, input CreateTransactionInput) (string, error) {
	ctx, span := c.otel.Start(ctx, "CreateTransaction")
	defer span.End()
	if !c.transactionAuthorizer.IsTransactionAllowed(ctx) {
		return "", errs.ErrTransactionNotAllowed
	}

	var transactionID string

	err := c.userRepository.UpdateBalance(ctx, input.SenderID.String(), input.ReceiverID.String(), func(sender, receiver *entity.User) (*entity.Transaction, error) {
		if sender.IsMerchant() {
			return nil, errs.ErrMerchantCannotSendMoney
		}

		err := sender.Withdraw(input.Amount)
		if err != nil {
			return nil, err
		}

		err = receiver.Deposit(input.Amount)
		if err != nil {
			return nil, err
		}

		transaction, err := entity.NewTransaction(input.Amount, sender.ID(), receiver.ID())
		if err != nil {
			return nil, err
		}

		transactionID = transaction.ID()

		eventTransaction := event.NewCreateTransactionEventV1(transactionID, input.Amount, input.SenderID, input.ReceiverID)
		err = c.queue.Send(ctx, eventTransaction.ToJSON())
		if err != nil {
			return nil, err
		}

		return transaction, nil
	})
	if err != nil {
		return "", err
	}

	return transactionID, err
}

func NewCreateTransaction(
	userRepository UserRepository,
	transactionAuthorizer TransactionAuthorizerGateway,
	queue Queue,
	otel telemetry.Telemetry,
) *CreateTransaction {
	return &CreateTransaction{
		userRepository:        userRepository,
		transactionAuthorizer: transactionAuthorizer,
		queue:                 queue,
		otel:                  otel,
	}
}
