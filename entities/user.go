package entities

import (
	"Shopee_UMS/utils"
)

type Account struct {
	Username string
	Password string
}

type User struct {
	Account    *Account
	Nickname   string
	ProfileUri string
}

func ValidateUserNickname(name string) error {
	if len(name) > 255 {
		return utils.ValidationError{"nickname is too long"}
	}
	if len(name) == 0 {
		return utils.ValidationError{"nickname is empty"}
	}
	return nil
}
