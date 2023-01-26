package jwt_token_helper

import (
	"crypto"
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var METHOD_NAME = "EdDSA"

type myAuthClaims struct {
	name string
	jwt.RegisteredClaims
}

type TokenHelper struct {
	secret      interface{}
	pub         interface{}
	expiredTime time.Duration
}

func New(secret []byte, expired time.Duration) (*TokenHelper, error) {
	key, err := jwt.ParseEdPrivateKeyFromPEM(secret)
	if err != nil {
		return nil, err
	}

	var pub crypto.PublicKey
	if pkey, ok := key.(ed25519.PrivateKey); ok {
		pub = pkey.Public()
	} else {
		return nil, errors.New("create public key failed")
	}

	return &TokenHelper{
		secret:      key,
		expiredTime: expired,
		pub:         pub,
	}, nil
}

func (t *TokenHelper) CreateToken(name string) (string, int64, error) {
	expired := time.Now().Add(t.expiredTime)

	claims := myAuthClaims{
		name: "aa",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expired),
			Issuer:    name,
			Subject:   "login",
		},
	}

	j := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token, err := j.SignedString(t.secret)
	if err != nil {
		return "", 0, err
	}

	return token, expired.Unix(), nil
}

func (t *TokenHelper) VerifyToken(jwttoken string) (bool, string) {
	token, err := verifyToken(t.pub, jwttoken)
	if err != nil {
		return false, err.Error()
	}

	if !token.Valid {
		return false, ""
	}

	return true, token.Claims.(*myAuthClaims).RegisteredClaims.Issuer
}

func verifyToken(publicKey interface{}, token string) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(token, &myAuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*myAuthClaims); !ok {
		return nil, fmt.Errorf("unexpected claims: %v", t.Claims)
	} else if claims == nil ||
		claims.RegisteredClaims.ExpiresAt.Before(time.Now()) ||
		claims.RegisteredClaims.Issuer == "" {
		return nil, fmt.Errorf("invalid token")
	}
	return t, nil
}

// TODO verify
