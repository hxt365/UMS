package reposistory

import (
	"Shopee_UMS/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccounts_Get(t *testing.T) {
	setup()
	defer tearDown()
	repo := NewAccounts(db)
	acc, err := repo.Get("user1")
	assert.Nil(t, err)
	assert.Equal(t, 1, acc.Id)
	assert.Equal(t, "user1", acc.Username)
	assert.True(t, entities.VerifyPassword(acc.Password, "secret"))
}
