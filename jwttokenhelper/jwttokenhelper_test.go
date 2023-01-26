package jwt_token_helper

import (
	"os"
	"strconv"
	"testing"
	"time"
)

var helper TokenHelper
var expired = time.Hour

func TestMain(m *testing.M) {
	file, err := os.ReadFile(".secret/JWT_KEY.pem")
	if err != nil {
		panic(err)
	}

	h, err := New([]byte(file), expired)
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

		if time.Now().Add(expired).Unix() < e {
			t.Error("ee")
		}

		r, r_name := helper.VerifyToken(token)
		if !r {
			t.Errorf("verify token failed, %v", r_name)
		}

		if r_name != name {
			t.Errorf("exprect name %s, but %s", name, r_name)
		}
	})
}

func TestTokenHelper_VerifyToken(t *testing.T) {

}

func Test_verifyToken(t *testing.T) {

}
