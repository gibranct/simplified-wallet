package usecase

import (
	"context"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UpdateBalance(ctx context.Context, senderID, receiverID string, updateFn func(sender, receiver *entity.User) (*entity.Transaction, error)) error
}

type TransactionAuthorizerGateway interface {
	IsTransactionAllowed(ctx context.Context) bool
}

type CreateTransaction struct {
	userRepository        UserRepository
	transactionAuthorizer TransactionAuthorizerGateway
}

type CreateTransactionInput struct {
	Amount     float64
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
}

func (c *CreateTransaction) Execute(ctx context.Context, input CreateTransactionInput) (string, error) {
	if !c.transactionAuthorizer.IsTransactionAllowed(ctx) {
		return "", errs.ErrTransactionNotAllowed
	}

	var transactionID string

	err := c.userRepository.UpdateBalance(ctx, input.SenderID.String(), input.ReceiverID.String(), func(sender, receiver *entity.User) (*entity.Transaction, error) {
		err := sender.Withdraw(input.Amount)
		if err != nil {
			return nil, err
		}

		receiver.Deposit(input.Amount)

		transaction := entity.NewTransaction(input.Amount, sender.ID(), receiver.ID())

		transactionID = transaction.ID()

		return transaction, nil
	})

	return transactionID, err
}

func NewCreateTransaction(
	userRepository UserRepository,
	transactionAuthorizer TransactionAuthorizerGateway,
) CreateTransaction {
	return CreateTransaction{
		userRepository:        userRepository,
		transactionAuthorizer: transactionAuthorizer,
	}
}
