package vo

import (
	"errors"
	"regexp"
)

var ErrInvalidEmail = errors.New("invalid email")

type Email struct {
	value string
}

func NewEmail(value string) (*Email, error) {
	matchEmail := regexp.MustCompile("^(.+)@(.+)$").MatchString(value)
	if !matchEmail {
		return nil, ErrInvalidEmail
	}
	return &Email{
		value: value,
	}, nil
}

func (e *Email) GetValue() string {
	return e.value
}
