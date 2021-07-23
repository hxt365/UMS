package api

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/usecases"
	"Shopee_UMS/utils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type authUsecaseStub struct{}

func (au *authUsecaseStub) Authenticate(username, password string) (int, error) {
	return 1, nil
}

func TestAuth_TokenLoginSuccessfully(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})

	w := login(t, s)

	authCookie := extractAuthCookie(w)
	assert.NotNil(t, authCookie, "should set auth cookie after login")
	token := authCookie.Value
	jwtAuth, _ := s.auth.(*entities.JwtAuthenticator)
	claim, err := jwtAuth.ValidateToken(token)
	assert.Nil(t, err, "invalid jwt token")
	jwtClaim := claim.(*entities.JwtClaim)
	assert.Equal(t, 1, jwtClaim.Uid)

	var user usecases.UserData
	err = json.Unmarshal(w.Body.Bytes(), &user)
	assert.Nil(t, err, "could not decode response body")
	assert.Equal(t, "user", user.Username)
	assert.Equal(t, "nickname", user.Nickname)
	assert.Equal(t, "s3://something.com", user.ProfilePicUri)
}

type authUsecaseWrongStub struct{}

func (au *authUsecaseWrongStub) Authenticate(username, password string) (int, error) {
	return 0, utils.AuthError{"wrong password"}
}

func TestAuth_TokenLoginFailDueToWrongPassword(t *testing.T) {
	s := NewTestServer(&authUsecaseWrongStub{}, nil)

	reqJson := `{"username": "user", "password": "secret"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(reqJson))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuth_Logout(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, &userUsecaseStub{})

	w := login(t, s)
	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	req.AddCookie(extractAuthCookie(w))
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	authCookie := extractAuthCookie(w)
	assert.Equal(t, "", authCookie.Value)
	assert.True(t, authCookie.Expires.Before(time.Now().UTC()))
}

func TestAuth_LogoutFailDueToUnauthenticated(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, nil)

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
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
		if c.Name == "auth-token" {
			return c
		}
	}
	return nil
}
