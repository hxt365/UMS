package usecases

import (
	"Shopee_UMS/entities"
	"fmt"
)

type AuthUsecaser interface {
	Authenticate(username, password string) error
}

type authUsecase struct {
	accounts AccountRepository
}

func (au *authUsecase) Authenticate(username, password string) error {
	hash := au.accounts.GetHashPassword(username)
	if !entities.CheckPasswordHash(hash, password) {
		return fmt.Errorf("wrong username or password")
	}
	return nil
}
