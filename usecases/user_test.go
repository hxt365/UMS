package usecases

import (
	"Shopee_UMS/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

type userRepoStub struct{}

func (ur *userRepoStub) GetByUsername(username string) (*entities.User, error) {
	return &entities.User{
		Account: &entities.Account{
			Username: username,
			Password: "password",
		},
		Nickname:   "nickname",
		ProfileUri: "s3://something.com",
	}, nil
}

func (ur *userRepoStub) UpdateNickname(username, nickname string) error {
	return nil
}

func (ur *userRepoStub) UpdateProfileUri(username, url string) error {
	return nil
}

type staticStorageStub struct{}

func (ss *staticStorageStub) UploadFile(data []byte, name string) (string, error) {
	return "s3://someuri.com", nil
}

func TestUserUsecase_GetInformation(t *testing.T) {
	uu := &userUsecase{users: &userRepoStub{}}
	user, err := uu.GetInformation("shopee")
	assert.Nil(t, err, "could not get user information")
	assert.Equal(t, user.Username, "shopee")
	assert.Equal(t, user.ProfileUri, "s3://something.com")
	assert.Equal(t, user.Nickname, "nickname")
}

func TestUserUsecase_ChangeNickname(t *testing.T) {
	uu := &userUsecase{users: &userRepoStub{}}
	err := uu.ChangeNickname("user1", "user2")
	assert.Nil(t, err, "could not change nickname")
}

func TestUserUsecase_UploadProfilePicture(t *testing.T) {
	uu := &userUsecase{
		users:   &userRepoStub{},
		storage: &staticStorageStub{},
	}
	err := uu.UploadProfilePicture("username", &Photo{})
	assert.Nil(t, err, "could not upload profile picture")
}
