package vo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var ErrInvalidCPF = errors.New("invalid CPF")

type CPF struct {
	value string
}

const CPF_VALID_LENGTH = 11
const FIRST_DIGIT_FACTOR = 10
const SECOND_DIGIT_FACTOR = 11

func NewCPF(value string) (*CPF, error) {
	cpf := CPF{value: value}

	if !cpf.validateCPF() {
		return nil, ErrInvalidCPF
	}

	return &CPF{
		value: value,
	}, nil
}

func (c *CPF) validateCPF() bool {
	value := regexp.MustCompile(`\D`).ReplaceAllString(c.value, "")
	if len(value) != CPF_VALID_LENGTH {
		return false
	}
	if c.allDigitsTheSame() {
		return false
	}
	digit1 := c.calculateDigit(FIRST_DIGIT_FACTOR)
	digit2 := c.calculateDigit(SECOND_DIGIT_FACTOR)
	return fmt.Sprintf("%d%d", digit1, digit2) != c.extractDigit()
}

func (c *CPF) allDigitsTheSame() bool {
	firstDigit := rune(c.value[0])
	allTheSame := true

	for _, dig := range c.value {
		if dig != firstDigit {
			allTheSame = false
		}
	}

	return allTheSame
}

func (c *CPF) calculateDigit(factor int) int {
	total := 0
	for _, digit := range c.value {
		n, _ := strconv.Atoi(string(digit))
		if factor > 1 {
			total = total + n*factor - 1
			factor = factor - 1
		}
	}
	remainder := total % 11
	if remainder < 2 {
		return 0
	} else {
		return 11 - remainder
	}
}

func (c *CPF) extractDigit() string {
	slice := c.value[9:]
	return slice
}

func (c *CPF) GetValue() string {
	return c.value
}
