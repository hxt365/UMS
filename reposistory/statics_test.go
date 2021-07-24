package reposistory

import (
	"Shopee_UMS/storage"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestStatics_UploadFile(t *testing.T) {
	s3, err := storage.New()
	assert.Nil(t, err, "could not connect to AWS S3")
	bucketName := os.Getenv("S3_BUCKET_TEST_NAME")
	repo := NewStatics(s3, bucketName)

	file, err := os.Open("../assets/small.txt")
	assert.Nil(t, err, "could not load sample test file")
	filename := "test_" + fmt.Sprint(time.Now().UTC().Unix())

	uri, err := repo.UploadFile(file, filename)
	assert.Nil(t, err, "should upload successfully")
	assert.Equal(t, fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, s3.Region, filename), uri)
}
