package usecases

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/utils"
	"fmt"
	"mime/multipart"
	"time"
)

type UserUsecaser interface {
	GetData(uid int) (*UserData, error)
	ChangeNickname(uid int, nickname string) (string, error)
	UploadProfilePicture(uid int, photo *Photo) (string, error)
}

type userUsecase struct {
	users   UserRepository
	storage StaticStorage
}

type UserData struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	ProfilePicUri string `json:"profilePictureUri"`
}

func (uu *userUsecase) GetData(uid int) (*UserData, error) {
	u, err := uu.users.Get(uid)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uu *userUsecase) ChangeNickname(uid int, nickname string) (string, error) {
	if err := entities.ValidateUserNickname(nickname); err != nil {
		return "", err
	}
	if err := uu.users.UpdateNickname(uid, nickname); err != nil {
		return "", err
	}
	return nickname, nil
}

type Photo struct {
	Data multipart.File
	Size int64
}

const MaxProfilePictureSizeMB = 2

func (uu *userUsecase) UploadProfilePicture(uid int, photo *Photo) (string, error) {
	if photo.Size > MaxProfilePictureSizeMB*1000000 {
		return "", utils.ValidationError{
			Err: fmt.Sprintf("profile picture size should not exceed %d MB", MaxProfilePictureSizeMB)}
	}
	filename := fmt.Sprint(uid) + "_" + fmt.Sprint(time.Now().UTC().Unix())
	uri, err := uu.storage.UploadFile(photo.Data, filename)
	if err != nil {
		return "", err
	}
	err = uu.users.UpdateProfilePicUri(uid, uri)
	return uri, err
}
