package vo

import (
	"fmt"
)

const (
	CommonUserType   = "common"
	MerchantUserType = "merchant"
)

type UserType struct {
	value string
}

func NewUserType(value string) (*UserType, error) {
	validUserTypes := []string{CommonUserType, MerchantUserType}
	for _, validType := range validUserTypes {
		if value == validType {
			return &UserType{value: value}, nil
		}
	}
	return &UserType{}, fmt.Errorf("invalid user type: %s", value)
}

func (u UserType) IsMerchant() bool {
	return u.value == MerchantUserType
}

func (u UserType) IsCommon() bool {
	return u.value == CommonUserType
}

func (u UserType) Value() string {
	return u.value
}
