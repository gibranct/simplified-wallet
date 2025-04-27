package vo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	CNPJ_VALID_LENGTH        = 14
	CNPJ_FIRST_DIGIT_FACTOR  = 5
	CNPJ_SECOND_DIGIT_FACTOR = 6
	CNPJ_DIGIT_MODULE        = 11
)

var ErrInvalidCNPJ = errors.New("invalid CNPJ")

type CNPJ struct {
	value string
}

func NewCNPJ(value string) (*CNPJ, error) {
	cnpj := &CNPJ{value: value}
	if !cnpj.validateCNPJ() {
		return nil, ErrInvalidCNPJ
	}
	cnpj.value = regexp.MustCompile(`\D`).ReplaceAllString(cnpj.value, "")
	return cnpj, nil
}

func (c *CNPJ) GetValue() string {
	return c.value
}

func (c *CNPJ) validateCNPJ() bool {
	value := regexp.MustCompile(`\D`).ReplaceAllString(c.value, "")
	if len(value) != CNPJ_VALID_LENGTH {
		return false
	}
	if c.allDigitsTheSame() {
		return false
	}

	digit1 := c.calculateFirstDigit()
	digit2 := c.calculateSecondDigit()

	calculatedDigits := fmt.Sprintf("%d%d", digit1, digit2)
	actualDigits := c.extractDigits()

	return calculatedDigits == actualDigits
}

func (c *CNPJ) allDigitsTheSame() bool {
	value := regexp.MustCompile(`\D`).ReplaceAllString(c.value, "")
	firstDigit := value[0]
	for i := 1; i < len(value); i++ {
		if value[i] != firstDigit {
			return false
		}
	}
	return true
}

func (c *CNPJ) calculateFirstDigit() int {
	value := regexp.MustCompile(`\D`).ReplaceAllString(c.value, "")

	sum := 0
	factor := CNPJ_FIRST_DIGIT_FACTOR

	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(value[i]))
		sum += digit * factor

		factor--
		if factor < 2 {
			factor = 9
		}
	}

	remainder := sum % CNPJ_DIGIT_MODULE
	if remainder < 2 {
		return 0
	}
	return CNPJ_DIGIT_MODULE - remainder
}

func (c *CNPJ) calculateSecondDigit() int {
	value := regexp.MustCompile(`\D`).ReplaceAllString(c.value, "")

	sum := 0
	factor := CNPJ_SECOND_DIGIT_FACTOR

	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(value[i]))
		sum += digit * factor

		factor--
		if factor < 2 {
			factor = 9
		}
	}

	remainder := sum % CNPJ_DIGIT_MODULE
	if remainder < 2 {
		return 0
	}
	return CNPJ_DIGIT_MODULE - remainder
}

func (c *CNPJ) extractDigits() string {
	value := regexp.MustCompile(`\D`).ReplaceAllString(c.value, "")
	return value[12:14]
}
