package reposistory

import (
	"Shopee_UMS/storage"
	"Shopee_UMS/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestStatics_UploadFile(t *testing.T) {
	s3, err := storage.NewS3()
	assert.Nil(t, err, "could not connect to AWS S3")
	bucketName := utils.MustEnv("TEST_S3_BUCKET_NAME")
	repo := NewTestStatics(s3)

	file, err := os.Open("../assets/small.txt")
	assert.Nil(t, err, "could not load sample test file")
	filename := "test_" + fmt.Sprint(time.Now().UTC().Unix())

	uri, err := repo.UploadFile(file, filename)
	assert.Nil(t, err, "should upload successfully")
	assert.Equal(t, fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, s3.Region, filename), uri)
}
