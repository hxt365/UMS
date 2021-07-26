package storage

import (
	"Shopee_UMS/utils"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime/multipart"
)

type S3 struct {
	Region string
	ss     *session.Session
}

func NewS3() (*S3, error) {
	accessId := utils.MustEnv("AWS_ACCESS_KEY_ID")
	secretKey := utils.MustEnv("AWS_SECRET_ACCESS_KEY")
	region := utils.MustEnv("AWS_REGION")

	ss, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessId, secretKey, ""),
		},
	)
	if err != nil {
		return nil, err
	}

	return &S3{ss: ss, Region: region}, nil
}

func (s3 *S3) UploadFile(bucketName, acl, filename string, file multipart.File) (string, error) {
	uploader := s3manager.NewUploader(s3.ss)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		ACL:    aws.String(acl),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	uri := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, s3.Region, filename)
	return uri, nil
}
