package vo

import (
	"unicode/utf8"

	"github.com.br/gibranct/simplified-wallet/internal/domain/errs"
)

type Name struct {
	value string
}

func NewName(value string) (*Name, error) {
	strSize := utf8.RuneCountInString(value)
	if strSize < 3 || strSize > 50 {
		return nil, errs.ErrNameLength
	}
	return &Name{value: value}, nil
}

func (n *Name) Value() string {
	return n.value
}
