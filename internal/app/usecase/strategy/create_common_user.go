package strategy

import (
	"context"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
)

type CreateCommonUserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	ExistsByCPF(ctx context.Context, cpf string) (bool, error)
}

type CreateCommonUser struct {
	repository CreateCommonUserRepository
}

func (cuc *CreateCommonUser) UserType() string {
	return vo.CommonUserType
}

func (cuc *CreateCommonUser) Execute(ctx context.Context, input CreateUserStrategyInput) (string, error) {
	exists, err := cuc.repository.ExistsByCPF(ctx, input.Document)
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
		input.Document,
		"",
		vo.CommonUserType,
	)
	if err != nil {
		return "", err
	}
	err = cuc.repository.Save(ctx, user)
	if err != nil {
		return "", err
	}
	return user.ID(), nil
}

func NewCreateCommonUser(repository CreateCommonUserRepository) *CreateCommonUser {
	return &CreateCommonUser{
		repository: repository,
	}
}
