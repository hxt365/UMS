package entities

import (
	"crypto/rsa"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claim interface{}

type Authenticator interface {
	Authenticate(credentials ...string) (Claim, error)
}

type TokenAuthenticator interface {
	GenerateToken(claim Claim) (string, error)
	ValidateToken(signedToken string) (Claim, error)
}

type JwtAuthenticator struct {
	SecretKey  *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	ExpSeconds int64
	Issuer     string
}

func (j *JwtAuthenticator) Authenticate(credentials ...string) (Claim, error) {
	token := credentials[0]
	return j.ValidateToken(token)
}

type JwtClaim struct {
	Username string
	jwt.StandardClaims
}

func (j *JwtAuthenticator) GenerateToken(claim Claim) (string, error) {
	claimMap, ok := claim.(map[string]string)
	if !ok {
		return "", fmt.Errorf("wrong JWT claim format")
	}
	username, ok := claimMap["username"]
	if !ok {
		return "", fmt.Errorf("missing username in JWT claim")
	}

	jwtClaim := &JwtClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Duration(j.ExpSeconds) * time.Second).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaim)
	signedToken, err := token.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *JwtAuthenticator) ValidateToken(signedToken string) (Claim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return j.PublicKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*JwtClaim)
	if !ok {
		return nil, fmt.Errorf("could not parse JWT claim")
	}

	if claim.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("JWT token expired")
	}

	return claim, nil
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
