package vo

import "github.com.br/gibranct/simplified-wallet/internal/domain/errs"

type Money struct {
	value int64
}

func NewMoney(value float64) (*Money, error) {
	if value < 0 {
		return nil, errs.ErrZeroOrNegativeAmount
	}
	return &Money{value: int64(value * 100)}, nil
}

func (m Money) Value() int64 {
	return m.value
}

func (m Money) Subtract(amount float64) (*Money, error) {
	other, err := NewMoney(amount)
	if err != nil {
		return nil, err
	}
	if m.value < other.value {
		return nil, errs.ErrInsufficientBalance
	}
	return &Money{value: m.value - other.value}, nil
}

func (m Money) Add(amount float64) (*Money, error) {
	other, err := NewMoney(amount)
	if err != nil {
		return nil, err
	}
	return &Money{value: m.value + other.value}, nil
}
