package vo

import (
	"fmt"
)

const (
	commonUserType   = "common"
	merchantUserType = "merchant"
)

type UserType struct {
	value string
}

func NewUserType(value string) (*UserType, error) {
	validUserTypes := []string{commonUserType, merchantUserType}
	for _, validType := range validUserTypes {
		if value == validType {
			return &UserType{value: value}, nil
		}
	}
	return &UserType{}, fmt.Errorf("invalid user type: %s", value)
}

func (u UserType) IsMerchant() bool {
	return u.value == merchantUserType
}

func (u UserType) IsCommon() bool {
	return u.value == commonUserType
}

func (u UserType) Value() string {
	return u.value
}
