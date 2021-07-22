package api

import (
	"Shopee_UMS/entities"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type authUsecaseStub struct{}

func (au *authUsecaseStub) Authenticate(username, password string) error {
	return nil
}

func TestAuth_TokenLoginSuccessfully(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, nil)

	reqJson := `{"username": "user", "password": "secret"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqJson))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	authCookie := extractAuthCookie(w)
	assert.NotNilf(t, authCookie, "should set auth cookie after login")
	token := authCookie.Value
	jwtAuth, _ := s.auth.(*entities.JwtAuthenticator)
	claim, err := jwtAuth.ValidateToken(token)
	assert.Nil(t, err, "invalid jwt token")
	jwtClaim := claim.(*entities.JwtClaim)
	assert.Equal(t, "user", jwtClaim.Username)
}

type authUsecaseWrongStub struct{}

func (au *authUsecaseWrongStub) Authenticate(username, password string) error {
	return fmt.Errorf("fail")
}

func TestAuth_TokenLoginFailDueToWrongPassword(t *testing.T) {
	s := NewTestServer(&authUsecaseWrongStub{}, nil)

	reqJson := `{"username": "user", "password": "secret"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqJson))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuth_Logout(t *testing.T) {
	s := NewTestServer(&authUsecaseStub{}, nil)

	reqJson := `{"username": "user", "password": "secret"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqJson))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest("POST", "/api/logout", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	authCookie := extractAuthCookie(w)
	assert.Equal(t, "", authCookie.Value)
	assert.Less(t, authCookie.Expires.Unix(), time.Now().Local().Unix())
}

func extractAuthCookie(w *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range w.Result().Cookies() {
		if c.Name == "auth-token" {
			return c
		}
	}
	return nil
}
