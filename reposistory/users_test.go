package reposistory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsers_Get(t *testing.T) {
	prepareTestDB()
	repo := NewUsers(db)
	u, err := repo.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "nickname", u.Nickname)
	assert.Equal(t, "s3://somewhere.com", u.ProfilePicUri)
}

func TestUsers_UpdateNickname(t *testing.T) {
	prepareTestDB()
	repo := NewUsers(db)
	err := repo.UpdateNickname(1, "new nickname")
	assert.Nil(t, err)
	u, err := repo.Get(1)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "new nickname", u.Nickname)
	assert.Equal(t, "s3://somewhere.com", u.ProfilePicUri)
}

func TestUsers_UpdateProfilePicUri(t *testing.T) {
	prepareTestDB()
	repo := NewUsers(db)
	err := repo.UpdateProfilePicUri(1, "s3://newlocation")
	assert.Nil(t, err)
	u, err := repo.Get(1)
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "user1", u.Username)
	assert.Equal(t, "nickname", u.Nickname)
	assert.Equal(t, "s3://newlocation", u.ProfilePicUri)
}
