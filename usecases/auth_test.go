package usecases

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type accountRepoStub struct{}

func (ar *accountRepoStub) GetHashPassword(string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secretpassword"), bcrypt.MinCost)
	return string(hash)
}

func TestAuthUsecase_Authenticate(t *testing.T) {
	au := &authUsecase{accounts: &accountRepoStub{}}
	err := au.Authenticate("username", "secretpassword")
	assert.Nil(t, err, "wrong password")
}
