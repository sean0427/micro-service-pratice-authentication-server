package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock "github.com/sean0427/micro-service-pratice-auth-domain/mock"
	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

func FuzzCreateToken(f *testing.F) {
	f.Add("name_any", "pw", "token_any", int64(111232), 1, 1)
	f.Add("deaff", "fea", "awefafwefeawf", int64(112311), 1, 1)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", int64(11121), 1, 1)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", int64(11121), 1, 1)

	f.Fuzz(func(t *testing.T, name, pw, token string, expiredTime int64, authRun int, redisRun int) {
		expired := time.Unix(expiredTime, 0)

		ctrl := gomock.NewController(t)

		auth := mock.NewMockauthTool(ctrl)
		redis := mock.NewMockredisSvc(ctrl)

		auth.EXPECT().CreateToken(name).Return(token, expired, nil).Times(1)
		redis.EXPECT().Set(gomock.Any(), token, name, expired).Return(nil).Times(redisRun)

		got, err := createToken(context.Background(), name, auth, redis)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if got.Token != token {
			t.Errorf("want: %s, got: %s", token, got.Token)
		}

		if got.Name != name {
			t.Errorf("want: %s, got: %s", name, got.Name)
		}
	})
}

var testCreateToken_Error = []struct {
	name    string
	params  *model.LoginInfo
	authRun struct {
		times   int
		err     error
		expired time.Time
	}
	redisRun struct {
		times int
		err   error
	}
}{
	{
		name: "auth return error",
		params: &model.LoginInfo{
			Name:     "test",
			Password: "any",
		},
		authRun: struct {
			times   int
			err     error
			expired time.Time
		}{
			times: 1,
			err:   errors.New("auth error"),
		},
		redisRun: struct {
			times int
			err   error
		}{
			times: 0,
			err:   nil,
		},
	},
	{
		name: "redis return error",
		params: &model.LoginInfo{
			Name:     "test",
			Password: "any",
		},
		authRun: struct {
			times   int
			err     error
			expired time.Time
		}{
			times:   1,
			err:     nil,
			expired: time.Unix(93213, 0),
		},
		redisRun: struct {
			times int
			err   error
		}{
			times: 1,
			err:   errors.New("redis error"),
		},
	},
}

func TestCreateToken_Error(t *testing.T) {
	const token = "any"
	for _, c := range testCreateToken_Error {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			auth := mock.NewMockauthTool(ctrl)
			redis := mock.NewMockredisSvc(ctrl)

			auth.
				EXPECT().
				CreateToken(c.params.Name).
				Return(token, c.authRun.expired, c.authRun.err).
				Times(1)
			redis.
				EXPECT().
				Set(gomock.Any(), token, c.params.Name, c.authRun.expired).
				Return(c.redisRun.err).
				Times(c.redisRun.times)

			got, err := createToken(context.Background(), c.params.Name, auth, redis)
			if got != nil || err == nil {
				t.Fatal("should get error")
			}

			// shuld be ordering
			if c.redisRun.err != nil {
				if !errors.Is(err, c.redisRun.err) {
					t.Errorf("want: %v, got: %v", c.redisRun.err, err)
				}
				return
			}

			if c.authRun.err != nil {
				if !errors.Is(err, c.authRun.err) {
					t.Errorf("want: %v, got: %v", c.authRun.err, err)
				}
				return
			}
		})
	}
}

func FuzzVerify(f *testing.F) {
	f.Add("", "test", "feafefawefwa")
	f.Add("", "tesfeafewfjt", "feafefawefwa")
	f.Add("", "tesfeafewfjt", "feafefawefefefwa")
	f.Add("test", "test", "feafefawefwa")

	f.Fuzz(func(t *testing.T, name, returnedName, token string) {
		ctrl := gomock.NewController(t)
		auth := mock.NewMockauthTool(ctrl)
		redis := mock.NewMockredisSvc(ctrl)

		auth.EXPECT().
			VerifyToken(token).
			Return(true, returnedName, nil).
			Times(1)

		redis.EXPECT().
			Get(gomock.Any(), token).
			Return(returnedName, nil).
			Times(1)

		result, err := verifyToken(context.Background(),
			name, token,
			auth, redis,
		)

		if err != nil {
			t.Error("expected err to be nil, got:", err)
		}
		if !result {
			t.Errorf("want: %v, got: %v", true, result)
		}

	})
}

type authReturn struct {
	err     error
	success bool
}

var testVerify_Error_cases = []struct {
	name        string
	inputname   string
	authReturn  authReturn
	redisReturn struct {
		times       int
		err         error
		returnValue string
	}
}{
	{
		name:      "auth failed",
		inputname: "feaf",
		authReturn: authReturn{
			err:     errors.New("auth failed"),
			success: false,
		},
		// redis return not be used
	},
	{
		name:      "redis return err",
		inputname: "feaf",
		authReturn: authReturn{
			success: true,
		},
		redisReturn: struct {
			times       int
			err         error
			returnValue string
		}{
			times:       1,
			err:         errors.New("redis return error"),
			returnValue: "",
		},
	},
	{
		name:      "redis return name failed",
		inputname: "feaf",
		authReturn: authReturn{
			success: true,
		},
		redisReturn: struct {
			times       int
			err         error
			returnValue string
		}{
			times:       1,
			err:         nil,
			returnValue: "eror name",
		},
	},
}

func TestVerify_Error(t *testing.T) {
	const token = "test-token"

	for _, c := range testVerify_Error_cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			auth := mock.NewMockauthTool(ctrl)
			redis := mock.NewMockredisSvc(ctrl)

			auth.EXPECT().
				VerifyToken(token).
				Return(c.authReturn.success, c.inputname, c.authReturn.err).
				Times(1)

			redis.EXPECT().
				Get(gomock.Any(), token).
				Return(c.redisReturn.returnValue, c.redisReturn.err).
				Times(c.redisReturn.times)

			got, err := verifyToken(context.Background(),
				c.inputname, token,
				auth, redis,
			)

			if got {
				t.Errorf("always be false for verify")
			}

			if !c.authReturn.success {
				if !errors.Is(err, c.authReturn.err) {
					t.Errorf("want: %v, got: %v", c.authReturn.err, err)
				}
				return
			}
			if c.redisReturn.returnValue != c.inputname || c.redisReturn.err != nil {
				if err == nil {
					t.Error("expected error, got none")
				}
				return
			}
		})
	}
}
