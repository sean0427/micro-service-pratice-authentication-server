package jwt_token_helper

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var METHOD_NAME = jwt.SigningMethodES512.Name

type TokenHelper struct {
	secret      []byte
	expiredTime time.Duration
}

func New(secret []byte, expired time.Duration) *TokenHelper {
	return &TokenHelper{
		secret:      secret,
		expiredTime: expired,
	}
}

func (t *TokenHelper) CreateToken(account string) (string, int64, error) {
	expired := time.Now().Add(t.expiredTime).Unix()

	claims := jwt.MapClaims{
		"authorized": true,
		"account":    account,
		"expireed":   expired,
	}

	j := jwt.NewWithClaims(jwt.GetSigningMethod(METHOD_NAME), claims)
	token, err := j.SignedString(t.secret)
	if err != nil {
		return "", 0, err
	}

	return token, expired, nil
}

func (t *TokenHelper) VerifyToken(jwttoken string) (bool, string) {
	token, err := verifyToken(t.secret, jwttoken)
	if err != nil {
		return false, ""
	}

	if !token.Valid {
		return false, ""
	}

	return true, token.Claims.(jwt.MapClaims)["account"].(string)
}

func verifyToken(secret []byte, token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if claims, ok := token.Claims.(jwt.MapClaims); !ok {
			return nil, fmt.Errorf("unexpected claims: %v", token.Claims)
		} else if !claims["authorized"].(bool) ||
			claims["exipre"].(float64) < float64(time.Now().Unix()) ||
			claims["account"].(string) == "" {
			return nil, fmt.Errorf("invalid token")
		}
		return secret, nil
	})
}

// TODO verify
