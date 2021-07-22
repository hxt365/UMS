package usecases

import (
	"fmt"
	"time"
)

type UserUsecaser interface {
	GetInformation(username string) (*UserData, error)
	ChangeNickname(username, nickname string) error
	UploadProfilePicture(username string, photo *Photo) error
}

type userUsecase struct {
	users   UserRepository
	storage StaticStorage
}

type UserData struct {
	Username   string
	Nickname   string
	ProfileUri string
}

func (uu *userUsecase) GetInformation(username string) (*UserData, error) {
	u, err := uu.users.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return &UserData{
		Username:   u.Account.Username,
		Nickname:   u.Nickname,
		ProfileUri: u.ProfileUri,
	}, nil
}

func (uu *userUsecase) ChangeNickname(username, nickname string) error {
	err := uu.users.UpdateNickname(username, nickname)
	return err
}

type Photo struct {
	name string
	data []byte
	size int64
}

const MaxProfilePictureSizeMb = 5 // 5 MB

func (uu *userUsecase) UploadProfilePicture(username string, photo *Photo) error {
	if photo.size > MaxProfilePictureSizeMb*1000000000 {
		return fmt.Errorf("profile picture size should nod exceed %d MB", MaxProfilePictureSizeMb)
	}
	filename := username + fmt.Sprint(time.Now().Unix())
	uri, err := uu.storage.UploadFile(photo.data, filename)
	if err != nil {
		return err
	}
	err = uu.users.UpdateProfileUri(username, uri)
	return err
}
