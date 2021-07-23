package usecases

import (
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

type userRepoStub struct{}

func (ur *userRepoStub) Get(uid int) (*UserData, error) {
	return &UserData{
		Id:            uid,
		Username:      "user",
		Nickname:      "nickname",
		ProfilePicUri: "s3://something.com",
	}, nil
}

func (ur *userRepoStub) UpdateNickname(uid int, nickname string) error {
	return nil
}

func (ur *userRepoStub) UpdateProfilePicUri(uid int, url string) error {
	return nil
}

type staticStorageStub struct{}

func (ss *staticStorageStub) UploadFile(data multipart.File, name string) (string, error) {
	return "s3://someuri.com", nil
}

func TestUserUsecase_GetInformation(t *testing.T) {
	uu := &userUsecase{users: &userRepoStub{}}
	user, err := uu.GetData(1)
	assert.Nil(t, err, "could not get user information")
	assert.Equal(t, user.Username, "user")
	assert.Equal(t, user.ProfilePicUri, "s3://something.com")
	assert.Equal(t, user.Nickname, "nickname")
}

func TestUserUsecase_ChangeNickname(t *testing.T) {
	uu := &userUsecase{users: &userRepoStub{}}
	nickname, err := uu.ChangeNickname(1, "user")
	assert.Nil(t, err, "could not change nickname")
	assert.Equal(t, "user", nickname)
}

func TestUserUsecase_UploadProfilePicture(t *testing.T) {
	uu := &userUsecase{
		users:   &userRepoStub{},
		storage: &staticStorageStub{},
	}
	uri, err := uu.UploadProfilePicture(1, &Photo{})
	assert.Nil(t, err, "could not upload profile picture")
	assert.Equal(t, "s3://someuri.com", uri)
}
