package strategy

import (
	"context"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
)

type CreateMerchantUserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	ExistsByCNPJ(ctx context.Context, cpf string) (bool, error)
}

type CreateMerchantUser struct {
	repository CreateMerchantUserRepository
	otel       telemetry.Telemetry
}

func (cuc *CreateMerchantUser) UserType() string {
	return vo.MerchantUserType
}

func (cuc *CreateMerchantUser) Execute(ctx context.Context, input CreateUserStrategyInput) (string, error) {
	ctx, span := cuc.otel.Start(ctx, "CreateMerchantUser")
	defer span.End()

	exists, err := cuc.repository.ExistsByCNPJ(ctx, input.Document)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errs.ErrCNPJAlreadyRegistered
	}

	user, err := entity.NewUser(
		input.Name,
		input.Email,
		input.Password,
		"",
		input.Document,
		vo.MerchantUserType,
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

func NewCreateMerchantUser(repository CreateMerchantUserRepository, otel telemetry.Telemetry) *CreateMerchantUser {
	return &CreateMerchantUser{
		repository: repository,
		otel:       otel,
	}
}
