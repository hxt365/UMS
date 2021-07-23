package usecases

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/utils"
)

type AuthUsecaser interface {
	Authenticate(username, password string) (int, error)
}

type authUsecase struct {
	accounts AccountRepository
}

type AccountData struct {
	Id       int
	Username string
	Password string // hashed
}

// Authenticate returns user id if succeeded
func (au *authUsecase) Authenticate(username, password string) (int, error) {
	acc, err := au.accounts.Get(username)
	if err != nil {
		return 0, err
	}
	if !entities.VerifyPassword(acc.Password, password) {
		return 0, utils.AuthError{"wrong username or password"}
	}
	return acc.Id, nil
}
