package usecase

import (
	"context"

	"github.com.br/gibranct/simplified-wallet/internal/app/usecase/strategy"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
)

type CreateUserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type CreateUserStrategy interface {
	UserType() string
	Execute(ctx context.Context, input strategy.CreateUserStrategyInput) (string, error)
}

type CreateUser struct {
	userRepository CreateUserRepository
	strategies     []CreateUserStrategy
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	Document string
	UserType string
}

func (cus *CreateUser) Execute(ctx context.Context, input CreateUserInput) (string, error) {
	exists, err := cus.userRepository.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errs.ErrEmailAlreadyRegistered
	}
	var cuStrategy CreateUserStrategy
	for _, s := range cus.strategies {
		if s.UserType() == input.UserType {
			cuStrategy = s
		}
	}
	if cuStrategy == nil {
		return "", errs.ErrUserTypeNotFound
	}
	userID, err := cuStrategy.Execute(ctx, strategy.CreateUserStrategyInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Document: input.Document,
	})
	if err != nil {
		return "", err
	}
	return userID, nil
}

func NewCreateUser(userRepository CreateUserRepository, strategies []CreateUserStrategy) *CreateUser {
	return &CreateUser{
		userRepository: userRepository,
		strategies:     strategies,
	}
}
