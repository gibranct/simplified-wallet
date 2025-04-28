package entity

import (
	"time"

	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/google/uuid"
)

type User struct {
	id        uuid.UUID
	name      string
	email     *vo.Email
	password  *vo.Password
	balance   float64
	cpf       *vo.CPF
	cnpj      *vo.CNPJ
	userType  *vo.UserType
	active    bool
	createdAt time.Time
	updatedAt time.Time
}

func (u *User) ID() string {
	return u.id.String()
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email.GetValue()
}

func (u *User) Password() string {
	return u.password.Value
}

func (u *User) Balance() float64 {
	return u.balance
}

func (u *User) CPF() string {
	if u.cpf == nil {
		return ""
	}
	return u.cpf.GetValue()
}

func (u *User) CNPJ() string {
	if u.cnpj == nil {
		return ""
	}
	return u.cnpj.GetValue()
}

func (u *User) UserType() string {
	return u.userType.Value()
}

func (u *User) Active() bool {
	return u.active
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func NewUser(name, email, password, cpf, cnpj string, userType string) (*User, error) {
	id := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	userTypeEnum, err := vo.NewUserType(userType)
	if err != nil {
		return nil, err
	}

	if userTypeEnum.IsMerchant() && cnpj == "" {
		return nil, errs.ErrCNPJMustBeProvidedForMerchant
	}

	if userTypeEnum.IsMerchant() && cpf != "" {
		return nil, errs.ErrMerchantCannotHaveCPF
	}

	if userTypeEnum.IsCommon() && cpf == "" {
		return nil, errs.ErrCPFMustBeProvidedForCommonUser
	}

	if userTypeEnum.IsCommon() && cnpj != "" {
		return nil, errs.ErrCommonCannotHaveCNPJ
	}

	cpfObj, err := vo.NewCPF(cpf)
	if err != nil && userTypeEnum.IsCommon() {
		return nil, err
	}

	cnpjObj, err := vo.NewCNPJ(cnpj)
	if err != nil && userTypeEnum.IsMerchant() {
		return nil, err
	}

	emailObj, err := vo.NewEmail(email)
	if err != nil {
		return nil, err
	}
	passwordObj, err := vo.NewPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		id:        id,
		name:      name,
		email:     emailObj,
		password:  passwordObj,
		balance:   0.0,
		cpf:       cpfObj,
		cnpj:      cnpjObj,
		userType:  userTypeEnum,
		active:    true,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	return &user, nil
}

// Deposit adds money to the user's balance
func (u *User) Deposit(amount float64) {
	if amount > 0 {
		u.balance += amount
	}
}

// Withdraw removes money from the user's balance
func (u *User) Withdraw(amount float64) bool {
	if amount > 0 && amount <= u.balance {
		u.balance -= amount
		return true
	}
	return false
}
