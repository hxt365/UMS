package reposistory

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/usecases"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccounts_Get(t *testing.T) {
	setup()
	defer tearDown()
	repo := NewAccounts(db, c)
	acc, err := repo.Get("user1")
	assert.Nil(t, err)
	assert.Equal(t, 1, acc.Id)
	assert.Equal(t, "user1", acc.Username)
	assert.True(t, entities.VerifyPassword(acc.Password, "secret"))
}

func TestAccounts_GetWithRealCache(t *testing.T) {
	setup()
	defer tearDown()
	repo := NewAccounts(db, c)
	acc, err := repo.Get("user1")
	assert.Nil(t, err)
	assert.Equal(t, 1, acc.Id)
	assert.Equal(t, "user1", acc.Username)
	assert.True(t, entities.VerifyPassword(acc.Password, "secret"))

	var ad usecases.AccountData
	err = repo.c.LoadJson("accounts:user1", &ad)
	assert.Nil(t, err)
	assert.Equal(t, 1, acc.Id)
	assert.Equal(t, "user1", acc.Username)
	assert.True(t, entities.VerifyPassword(acc.Password, "secret"))

	err = repo.c.StoreJson("accounts:user1", &usecases.AccountData{
		Username: "from cache",
	})
	assert.Nil(t, err)

	acc, err = repo.Get("user1")
	assert.Nil(t, err)
	assert.Equal(t, "from cache", acc.Username)
}

func TestAccounts_GetWithMockCache(t *testing.T) {
	setup()
	defer tearDown()
	cache := &RedisMock{}
	repo := NewAccounts(db, cache)

	acc, err := repo.Get("user1")
	assert.Nil(t, err)
	assert.Equal(t, 1, acc.Id)
	assert.Equal(t, "user1", acc.Username)
	assert.True(t, entities.VerifyPassword(acc.Password, "secret"))
	assert.True(t, cache.CalledLoad())
	assert.True(t, cache.CalledStore())
	assert.False(t, cache.CalledDelete())

	_, err = repo.Get("user1")
	assert.Nil(t, err)
	assert.True(t, cache.CalledLoadTwice())
	assert.False(t, cache.CalledStoreTwice())
}
