package vo

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooShort = errors.New("password must be at least 6 characters long")

type Password struct {
	Value string
}

func (p *Password) Compare(hashedPassword, value string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(value))
	return err == nil
}

func NewPassword(value string) (*Password, error) {
	if utf8.RuneCountInString(value) < 6 {
		return nil, ErrPasswordTooShort
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Password{
		Value: string(hashedPassword),
	}, nil
}
