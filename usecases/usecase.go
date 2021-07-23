package usecases

import "mime/multipart"

type Usecases struct {
	Auth AuthUsecaser
	User UserUsecaser
}

func New(ar AccountRepository, ur UserRepository, ss StaticStorage) *Usecases {
	return &Usecases{
		Auth: &authUsecase{accounts: ar},
		User: &userUsecase{
			users:   ur,
			storage: ss,
		},
	}
}

type UserRepository interface {
	Get(uid int) (*UserData, error)
	UpdateNickname(uid int, nickname string) error
	UpdateProfilePicUri(uid int, url string) error
}

type AccountRepository interface {
	Get(username string) (*AccountData, error)
}

type StaticStorage interface {
	UploadFile(data multipart.File, name string) (string, error)
}
