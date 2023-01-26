package jwt_token_helper

import (
	"crypto"
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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

func (t *TokenHelper) CreateToken(name string) (string, time.Time, error) {
	expired := time.Now().Add(t.expiredTime)

	claims := myAuthClaims{
		name: "My Token",
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
		return "", time.Time{}, err
	}

	return token, expired, nil
}

func (t *TokenHelper) VerifyToken(tokenStr string) (bool, string, error) {
	token, err := verifyToken(t.pub, tokenStr)
	if err != nil {
		return false, "", err
	}

	if !token.Valid {
		return false, "", nil
	}

	return true, token.Claims.(*myAuthClaims).RegisteredClaims.Issuer, nil
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
