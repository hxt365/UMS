package entities

import (
	"Shopee_UMS/utils"
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
		if VerifyPassword(string(hash), tt.pwd) != tt.expect {
			assert.Equal(t, tt.expect, !tt.expect, "wrong hash password")
		}
	}
}

func TestJwtAuthenticator_GenerateToken(t *testing.T) {
	j := &JwtAuthenticator{
		PrivateKey: utils.SampleRSAPrivateKey(),
		PublicKey:  utils.SampleRSAPublicKey(),
		ExpSeconds: 60,
		Issuer:     "UMS",
	}
	claim := map[string]interface{}{
		"uid":  1,
		"csrf": "sometoken",
	}
	token, err := j.GenerateToken(claim)
	assert.Nil(t, err, "could not generate token")

	extractClaim, err := j.ValidateToken(token)
	jwtClaim, ok := extractClaim.(*JwtClaim)
	assert.True(t, ok, "could not convert extracted claim to JwtClaim")
	assert.Equal(t, jwtClaim.Uid, 1)
	assert.Equal(t, jwtClaim.Issuer, "UMS")
}
