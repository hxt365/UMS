package reposistory

import (
	"Shopee_UMS/usecases"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsers_Get(t *testing.T) {
	setup()
	defer tearDown()
	repo := NewUsers(db, c)
	u, err := repo.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "nickname", u.Nickname)
	assert.Equal(t, "s3://somewhere.com", u.ProfilePicUri)
}

func TestUsers_GetWithRealCache(t *testing.T) {
	setup()
	defer tearDown()
	repo := NewUsers(db, c)
	u, err := repo.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "nickname", u.Nickname)
	assert.Equal(t, "s3://somewhere.com", u.ProfilePicUri)

	var ud usecases.UserData
	err = repo.c.LoadJson("users:1", &ud)
	assert.Nil(t, err)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "nickname", u.Nickname)
	assert.Equal(t, "s3://somewhere.com", u.ProfilePicUri)

	err = c.StoreJson("users:1", &usecases.UserData{
		Username: "from cache",
	})
	assert.Nil(t, err)

	u, err = repo.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, "from cache", u.Username)
}

func TestUsers_GetWithMockCache(t *testing.T) {
	setup()
	defer tearDown()
	cache := &RedisMock{}
	repo := NewUsers(db, cache)

	u, err := repo.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "nickname", u.Nickname)
	assert.Equal(t, "s3://somewhere.com", u.ProfilePicUri)
	assert.True(t, cache.CalledLoad())
	assert.True(t, cache.CalledStore())

	_, err = repo.Get(1)
	assert.Nil(t, err)
	assert.True(t, cache.CalledLoadTwice())
	assert.False(t, cache.CalledStoreTwice())
}

func TestUsers_UpdateNickname(t *testing.T) {
	setup()
	defer tearDown()
	repo := NewUsers(db, c)
	err := repo.UpdateNickname(1, "new nickname")
	assert.Nil(t, err)
	u, err := repo.Get(1)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "new nickname", u.Nickname)
	assert.Equal(t, "s3://somewhere.com", u.ProfilePicUri)
}

func TestUsers_UpdateNicknameWithMockCache(t *testing.T) {
	setup()
	defer tearDown()
	cache := &RedisMock{}
	repo := NewUsers(db, cache)

	err := repo.UpdateNickname(1, "new nickname")
	assert.Nil(t, err)
	assert.True(t, cache.CalledDelete())
	assert.False(t, cache.CalledLoad())
	assert.False(t, cache.CalledStore())

	u, err := repo.Get(1)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "new nickname", u.Nickname)
	assert.Equal(t, "s3://somewhere.com", u.ProfilePicUri)
	assert.True(t, cache.CalledLoad())
	assert.True(t, cache.CalledStore())
}

func TestUsers_UpdateProfilePicUri(t *testing.T) {
	setup()
	defer tearDown()
	repo := NewUsers(db, c)
	err := repo.UpdateProfilePicUri(1, "s3://newlocation")
	assert.Nil(t, err)
	u, err := repo.Get(1)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "nickname", u.Nickname)
	assert.Equal(t, "s3://newlocation", u.ProfilePicUri)
}

func TestUsers_UpdateProfilePicUriWithMockCache(t *testing.T) {
	setup()
	defer tearDown()
	cache := &RedisMock{}
	repo := NewUsers(db, cache)

	err := repo.UpdateProfilePicUri(1, "s3://newlocation")
	assert.Nil(t, err)
	assert.True(t, cache.CalledDelete())
	assert.False(t, cache.CalledLoad())
	assert.False(t, cache.CalledStore())

	u, err := repo.Get(1)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "nickname", u.Nickname)
	assert.Equal(t, "s3://newlocation", u.ProfilePicUri)
	assert.True(t, cache.CalledLoad())
	assert.True(t, cache.CalledStore())
}
