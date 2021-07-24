package reposistory

import (
	"mime/multipart"
)

type Storage interface {
	UploadFile(bucketName, acl, filename string, file multipart.File) (string, error)
}

type Statics struct {
	storage    Storage
	bucketName string
}

func NewStatics(storage Storage, bucketName string) *Statics {
	return &Statics{
		storage:    storage,
		bucketName: bucketName,
	}
}

func (s *Statics) UploadFile(file multipart.File, filename string) (string, error) {
	return s.storage.UploadFile(s.bucketName, "public-read", filename, file)
}
