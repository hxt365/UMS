package entities

import (
	"Shopee_UMS/test_utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthenticator_CheckPasswordHash(t *testing.T) {
	tests := []struct {
		pwd    string
		expect bool
	}{
		{"secret", true},
		{"secret", false},
	}

	for _, tt := range tests {
		hash := []byte("wrong hash")
		if tt.expect {
			hash, _ = bcrypt.GenerateFromPassword([]byte(tt.pwd), bcrypt.MinCost)
		}
		if CheckPasswordHash(string(hash), tt.pwd) != tt.expect {
			assert.Equal(t, tt.expect, !tt.expect, "wrong hash password")
		}
	}
}

func TestJwtAuthenticator_GenerateToken(t *testing.T) {
	j := &JwtAuthenticator{
		SecretKey:  test_utils.GenRSAPrivateKey(),
		PublicKey:  test_utils.GenRSAPublicKey(),
		ExpSeconds: 60,
		Issuer:     "UMS",
	}
	claim := map[string]string{
		"username": "John Doe",
	}
	token, err := j.GenerateToken(claim)
	assert.Nil(t, err, "could not generate token")

	extractClaim, err := j.ValidateToken(token)
	jwtClaim, ok := extractClaim.(*JwtClaim)
	assert.True(t, ok, "could not convert extracted claim to JwtClaim")
	assert.Equal(t, jwtClaim.Username, "John Doe")
	assert.Equal(t, jwtClaim.Issuer, "UMS")
}
