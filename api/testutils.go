package api

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/usecases"
	"Shopee_UMS/utils"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func NewTestServer(au usecases.AuthUsecaser, uu usecases.UserUsecaser) *Server {
	jwtAuth := entities.JwtAuthenticator{
		PrivateKey: utils.SampleRSAPrivateKey(),
		PublicKey:  utils.SampleRSAPublicKey(),
		ExpSeconds: 10,
		Issuer:     "Test Server",
	}
	s := &Server{
		mux: http.NewServeMux(),
		u: &usecases.Usecases{
			Auth: au,
			User: uu,
		},
		auth: &jwtAuth,
	}
	s.routes()
	return s
}

func login(t *testing.T, s *Server) *httptest.ResponseRecorder {
	reqJson := `{"username": "user", "password": "secret"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(reqJson))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	return w
}

func extractAuthCookie(w *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range w.Result().Cookies() {
		if c.Name == utils.AuthCookieKey {
			return c
		}
	}
	return nil
}

func loadFileBody(t *testing.T, path string) (*bytes.Buffer, string) {
	file, err := os.Open(path)
	assert.Nil(t, err, "could not open test photo")
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, err := writer.CreateFormFile("profilePicture", filepath.Base(path))
	assert.Nil(t, err, "could not create form file")
	_, err = io.Copy(part, file)
	assert.Nil(t, err, "could not write file to part")

	return body, writer.FormDataContentType()
}
