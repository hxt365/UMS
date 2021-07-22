package usecases

import "Shopee_UMS/entities"

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
	GetByUsername(username string) (*entities.User, error)
	UpdateNickname(username, nickname string) error
	UpdateProfileUri(username, url string) error
}

type AccountRepository interface {
	GetHashPassword(username string) string
}

type StaticStorage interface {
	UploadFile(data []byte, name string) (string, error)
}
