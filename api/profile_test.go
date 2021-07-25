package api

import (
	"Shopee_UMS/usecases"
	"Shopee_UMS/utils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type userUsecaseStub struct{}

func (us *userUsecaseStub) GetData(uid int) (*usecases.UserData, error) {
	return &usecases.UserData{
		Id:            uid,
		Username:      "user",
		Nickname:      "nickname",
		ProfilePicUri: "s3://something.com",
	}, nil
}
func (us *userUsecaseStub) ChangeNickname(uid int, nickname string) (string, error) {
	return nickname, nil
}
func (us *userUsecaseStub) UploadProfilePicture(uid int, photo *usecases.Photo) (string, error) {
	if photo.Size > usecases.MaxProfilePictureSizeMB*1000000 {
		return "", utils.ValidationError{"too large photo"}
	}
	return "s3://somewhere.com", nil
}

func TestGetProfile(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})

	w := login(t, s)
	req := httptest.NewRequest("GET", "/api/profile", nil)
	req.AddCookie(extractAuthCookie(w))
	req.Header.Add("X-CSRFToken", w.Header().Get("X-CSRFToken"))
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var user usecases.UserData
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.Nil(t, err, "could not decode response body")
	assert.Equal(t, "user", user.Username)
	assert.Equal(t, "nickname", user.Nickname)
	assert.Equal(t, "s3://something.com", user.ProfilePicUri)
}

func TestHandleProfileFailDueToUnauthorized(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})

	tests := []struct {
		method string
		expect int
	}{
		{"GET", http.StatusUnauthorized},
		{"PUT", http.StatusUnauthorized},
		{"POST", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, "/api/profile", nil)
		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)
		assert.Equal(t, tt.expect, w.Code)
	}
}

func TestChangeNickName(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})

	w := login(t, s)
	reqJson := `{"nickname": "abcdef"}`
	req := httptest.NewRequest("PUT", "/api/profile", strings.NewReader(reqJson))
	req.AddCookie(extractAuthCookie(w))
	req.Header.Add("X-CSRFToken", w.Header().Get("X-CSRFToken"))
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var user usecases.UserData
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.Nil(t, err, "could not decode response body")
	assert.Equal(t, "user", user.Username)
	assert.Equal(t, "abcdef", user.Nickname)
	assert.Equal(t, "s3://something.com", user.ProfilePicUri)
}

func TestChangeNickNameFailDueToCSRF(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})

	w := login(t, s)
	reqJson := `{"nickname": "abcdef"}`
	req := httptest.NewRequest("PUT", "/api/profile", strings.NewReader(reqJson))
	req.AddCookie(extractAuthCookie(w))
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUploadProfilePicture(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})
	photoPath := "./../assets/test_photo.png"

	w := login(t, s)
	body, contentType := loadFileBody(t, photoPath)
	req := httptest.NewRequest("POST", "/api/profile-picture", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("X-CSRFToken", w.Header().Get("X-CSRFToken"))
	req.AddCookie(extractAuthCookie(w))
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var user usecases.UserData
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.Nil(t, err, "could not decode response body")
	assert.Equal(t, "user", user.Username)
	assert.Equal(t, "nickname", user.Nickname)
	assert.Equal(t, "s3://somewhere.com", user.ProfilePicUri)
}

func TestUploadProfilePictureFailDueToTooLargePhoto(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})
	photoPath := "./../assets/test_large_photo.jpg"

	w := login(t, s)
	body, contentType := loadFileBody(t, photoPath)
	req := httptest.NewRequest("POST", "/api/profile-picture", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("X-CSRFToken", w.Header().Get("X-CSRFToken"))
	req.AddCookie(extractAuthCookie(w))
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadProfilePictureFailDueToUnsupportedFile(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})
	photoPath := "./../assets/test_file.md"

	w := login(t, s)
	body, contentType := loadFileBody(t, photoPath)
	req := httptest.NewRequest("POST", "/api/profile-picture", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("X-CSRFToken", w.Header().Get("X-CSRFToken"))
	req.AddCookie(extractAuthCookie(w))
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
