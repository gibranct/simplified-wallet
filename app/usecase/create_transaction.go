package usecase

import (
	"context"

	"github.com.br/gibranct/simplified-wallet/app/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdateBalance(ctx context.Context, senderID, receiverID string, updateFn func(sender, receiver *domain.User) (string, error)) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error)
}

type TransactionAuthorizerGateway interface {
	IsTransactionAllowed(ctx context.Context) bool
}

type CreateTransaction struct {
	userRepository        UserRepository
	transactionRepository TransactionRepository
	transactionAuthorizer TransactionAuthorizerGateway
}

type CreateTransactionInput struct {
	Amount     float64
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
}

func (c *CreateTransaction) Execute(ctx context.Context, input CreateTransactionInput) (string, error) {
	sender, err := c.userRepository.GetUserByID(ctx, input.SenderID)
	if err != nil {
		return "", domain.ErrSenderNotFound
	}

	receiver, err := c.userRepository.GetUserByID(ctx, input.ReceiverID)
	if err != nil {
		return "", domain.ErrReceiverNotFound
	}

	if !c.transactionAuthorizer.IsTransactionAllowed(ctx) {
		return "", domain.ErrTransactionNotAllowed
	}

	var transactionID string

	err = c.userRepository.UpdateBalance(ctx, sender.ID(), receiver.ID(), func(sender, receiver *domain.User) (string, error) {
		if !sender.Withdraw(input.Amount) {
			return "", domain.ErrNotEnoughMoney
		}
		receiver.Deposit(input.Amount)

		transaction := domain.NewTransaction(input.Amount, sender.ID(), receiver.ID())

		savedTransaction, err := c.transactionRepository.CreateTransaction(ctx, transaction)
		if err != nil {
			return "", err
		}

		transactionID = savedTransaction.ID()

		return savedTransaction.ID(), nil
	})

	return transactionID, err
}

func NewCreateTransaction(
	userRepository UserRepository,
	transactionRepository TransactionRepository,
	transactionAuthorizer TransactionAuthorizerGateway,
) CreateTransaction {
	return CreateTransaction{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		transactionAuthorizer: transactionAuthorizer,
	}
}
