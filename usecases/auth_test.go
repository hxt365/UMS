package usecases

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type accountRepoStub struct{}

func (ar *accountRepoStub) Get(username string) (*AccountData, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secretpassword"), bcrypt.MinCost)
	return &AccountData{
		Id:       1,
		Username: "user",
		Password: string(hash),
	}, nil
}

func TestAuthenticate(t *testing.T) {
	au := &authUsecase{accounts: &accountRepoStub{}}
	uid, err := au.Authenticate("username", "secretpassword")
	assert.Nil(t, err, "wrong password")
	assert.Equal(t, 1, uid)
}
