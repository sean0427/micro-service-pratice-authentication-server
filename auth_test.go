package auth_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/sean0427/micro-service-pratice-auth-domain"
	mock "github.com/sean0427/micro-service-pratice-auth-domain/mock"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

func TestMain(m *testing.M) {
	orgCreateToken := CreateToken
	code := m.Run()
	CreateToken = orgCreateToken
	os.Exit(code)
}

func FuzzAuthService_Login(f *testing.F) {
	f.Add("name_any", "pw", "token_any", true, int64(111232), 1, 1)
	f.Add("deaff", "fea", "awefafwefeawf", true, int64(112311), 1, 1)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", true, int64(11121), 1, 1)
	f.Add("deaff", "fea", "awefafwefeawf", false, int64(0), 0, 0)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", false, int64(0), 0, 0)

	f.Fuzz(func(t *testing.T, name, pw, token string, success bool, expiredTime int64, authRun int, redisRun int) {
		ctrl := gomock.NewController(t)
		userService := mock.NewMockuserService(ctrl)
		auth := mock.NewMockauthTool(ctrl)
		redis := mock.NewMockredis(ctrl)

		userService.EXPECT().Authenticate(gomock.Any(), name, pw).Return(success, nil).Times(1)
		auth.EXPECT().CreateToken(name).Return(token, expiredTime, nil).Times(authRun)
		redis.EXPECT().Set(gomock.Any(), token, name, expiredTime).Return(nil).Times(redisRun)

		s := New(userService, redis, auth)
		got, err := s.Login(context.Background(), &model.LoginInfo{Name: name, Password: pw})

		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if got == nil && !success {
			t.Logf("login failed")
			return
		}
		if got.Token != token {
			t.Errorf("want: %s, got: %s", token, got.Token)
		}

		if got.Name != name {
			t.Errorf("want: %s, got: %s", name, got.Name)
		}
	})
}

var testProductService_Login_Error = []struct {
	name             string
	params           *model.LoginInfo
	userSericeReturn struct {
		success bool
		err     error
	}
	authRun struct {
		times   int
		err     error
		expired int64
	}
	redisRun struct {
		times int
		err   error
	}
}{
	{
		name: "user service return error",
		params: &model.LoginInfo{
			Name:     "test",
			Password: "any",
		},
		userSericeReturn: struct {
			success bool
			err     error
		}{

			success: false,
			err:     errors.New("user domain error"),
		},
		authRun: struct {
			times   int
			err     error
			expired int64
		}{
			times: 0,
			err:   nil,
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
		name: "auth return error",
		params: &model.LoginInfo{
			Name:     "test",
			Password: "any",
		},
		userSericeReturn: struct {
			success bool
			err     error
		}{

			success: true,
			err:     nil,
		},
		authRun: struct {
			times   int
			err     error
			expired int64
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
		userSericeReturn: struct {
			success bool
			err     error
		}{

			success: true,
			err:     nil,
		},
		authRun: struct {
			times   int
			err     error
			expired int64
		}{
			times:   1,
			err:     nil,
			expired: 93213,
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

func TestAuthService_Login_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	userService := mock.NewMockuserService(ctrl)
	auth := mock.NewMockauthTool(ctrl)
	redis := mock.NewMockredis(ctrl)

	const token = "any"
	for _, c := range testProductService_Login_Error {
		t.Run(c.name, func(t *testing.T) {
			userService.
				EXPECT().
				Authenticate(gomock.Any(), c.params.Name, c.params.Password).
				Return(c.userSericeReturn.success, c.userSericeReturn.err).
				Times(1)
			auth.
				EXPECT().
				CreateToken(c.params.Name).
				Return(token, c.authRun.expired, c.authRun.err).
				Times(c.authRun.times)
			redis.
				EXPECT().
				Set(gomock.Any(), token, c.params.Name, c.authRun.expired).
				Return(c.redisRun.err).
				Times(c.redisRun.times)

			s := New(userService, redis, auth)
			got, err := s.Login(context.Background(), c.params)
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

			if !errors.Is(err, c.userSericeReturn.err) {
				t.Errorf("want: %v, got: %v", c.userSericeReturn.err, err)
			}
		})
	}
}
