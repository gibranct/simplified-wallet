package usecase

import (
	"context"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
)

type CreateUserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByCPF(ctx context.Context, cpf string) (bool, error)
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
	exists, err := cus.userRepository.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errs.ErrEmailAlreadyRegistered
	}
	exists, err = cus.userRepository.ExistsByCPF(ctx, input.CPF)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errs.ErrCPFAlreadyRegistered
	}
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
