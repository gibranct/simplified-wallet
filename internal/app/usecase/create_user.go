package usecase

import (
	"context"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
)

type CreateUserRepository interface {
	Save(ctx context.Context, user *entity.User) error
}

type CreateUser struct {
	userRepository CreateUserRepository
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	CPF      string
}

func (cus *CreateUser) Execute(ctx context.Context, input CreateUserInput) (string, error) {
	user, err := entity.NewUser(
		input.Name,
		input.Email,
		input.Password,
		input.CPF,
		"",
		vo.CommonUserType,
	)
	if err != nil {
		return "", err
	}
	err = cus.userRepository.Save(ctx, user)
	if err != nil {
		return "", err
	}
	return user.ID(), nil
}

func NewCreateUser(userRepository CreateUserRepository) *CreateUser {
	return &CreateUser{
		userRepository: userRepository,
	}
}
