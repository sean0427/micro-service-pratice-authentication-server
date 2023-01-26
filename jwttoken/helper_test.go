package jwt_token_helper

import (
	"strconv"
	"testing"
	"time"
)

var helper TokenHelper

// TODO: to make testing running oj
const key = "-----BEGIN PRIVATE KEY-----\nMC4CAQAwBQYDK2VwBCIEIH9yFS4LNJpBmLd5+TP9nj8w9j9O7j1sdCxRITcyGYnK\n-----END PRIVATE KEY-----"

var expired = time.Hour

func TestMain(m *testing.M) {
	h, err := New([]byte(key), expired)
	if err != nil {
		panic(err)
	}

	helper = *h

	m.Run()
}

func FuzzTokenHelper_CreateToken_And_VerifyToken(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(strconv.Itoa(i))
	}

	f.Fuzz(func(t *testing.T, name string) {
		token, e, err := helper.CreateToken(name)

		if err != nil {
			t.Error(err)
			return
		}

		if token == "" {
			t.Error("token is empty")
		}

		if time.Now().Add(expired).Before(e) {
			t.Error("ee")
		}

		r, r_name, err := helper.VerifyToken(token)
		if !r {
			t.Errorf("verify token failed, %v", r_name)
		}

		if r_name != name {
			t.Errorf("exprect name %s, but %s", name, r_name)
		}

		if err != nil {
			t.Errorf("expect err is nil, but %v", err)
		}
	})
}

func TestTokenHelper_VerifyToken(t *testing.T) {

}

func Test_verifyToken(t *testing.T) {

}
