package model

import (
	"database/sql"
	"time"

	"github.com.br/gibranct/simplified-wallet/internal/domain/entity"
	"github.com/google/uuid"
)

type UserModel struct {
	ID        string         `db:"id"`
	Name      string         `db:"name"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	Balance   int64          `db:"balance"`
	CPF       sql.NullString `db:"cpf"`
	CNPJ      sql.NullString `db:"cnpj"`
	UserType  string         `db:"user_type"`
	Active    bool           `db:"active"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func NewUserModelFrom(u *entity.User) *UserModel {
	return &UserModel{
		ID:       u.ID(),
		Name:     u.Name(),
		Email:    u.Email(),
		Password: u.Password(),
		Balance:  u.Balance(),
		CPF: sql.NullString{
			String: u.CPF(),
			Valid:  u.CPF() != "",
		},
		CNPJ: sql.NullString{
			String: u.CNPJ(),
			Valid:  u.CNPJ() != "",
		},
		UserType:  u.UserType(),
		Active:    u.Active(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}
}

func (um *UserModel) ToEntity() (*entity.User, error) {
	user, err := entity.CreateUser(
		uuid.MustParse(um.ID),
		float64(um.Balance/100),
		um.Name,
		um.Email,
		um.Password,
		um.CPF.String,
		um.CNPJ.String,
		um.UserType,
		um.CreatedAt,
		um.UpdatedAt,
		um.Active,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
