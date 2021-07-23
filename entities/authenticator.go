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
	Uid int
	jwt.StandardClaims
}

func (j *JwtAuthenticator) GenerateToken(claim Claim) (string, error) {
	claimMap, ok := claim.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("wrong JWT claim format")
	}
	_, ok = claimMap["uid"]
	if !ok {
		return "", fmt.Errorf("missing user id in JWT claim")
	}
	uid, ok := claimMap["uid"].(int)
	if !ok {
		return "", fmt.Errorf("user id in JWT claim should be int")
	}

	jwtClaim := &JwtClaim{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Duration(j.ExpSeconds) * time.Second).Unix(),
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

	exp := time.Unix(claim.ExpiresAt, 0)
	if exp.Before(time.Now().UTC()) {
		return nil, fmt.Errorf("JWT token expired")
	}

	return claim, nil
}

func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
