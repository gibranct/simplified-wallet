package vo

import (
	"fmt"
)

type UserType struct {
	value string
}

func NewUserType(value string) (*UserType, error) {
	validUserTypes := []string{"common", "merchant"}
	for _, validType := range validUserTypes {
		if value == validType {
			return &UserType{value: value}, nil
		}
	}
	return &UserType{}, fmt.Errorf("invalid user type: %s", value)
}

func (u UserType) Value() string {
	return u.value
}
