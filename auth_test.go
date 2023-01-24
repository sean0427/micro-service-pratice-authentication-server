package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/sean0427/micro-service-pratice-auth-domain"
	mock "github.com/sean0427/micro-service-pratice-auth-domain/mock"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

func FuzzAuthService_Login(f *testing.F) {
	f.Add("name_any", "pw", "token_any", true, 1, 1)
	f.Add("deaff", "fea", "awefafwefeawf", true, 1, 1)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", true, 1, 1)
	f.Add("deaff", "fea", "awefafwefeawf", false, 0, 0)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", false, 0, 0)

	f.Fuzz(func(t *testing.T, name, pw, token string, success bool, authRun int, redisRun int) {
		ctrl := gomock.NewController(t)
		userService := mock.NewMockuserService(ctrl)
		auth := mock.NewMockauthTool(ctrl)
		redis := mock.NewMockredis(ctrl)

		userService.EXPECT().Authenticate(gomock.Any(), name, pw).Return(success, nil).Times(1)
		auth.EXPECT().CreateToken(name).Return(token, nil).Times(authRun)
		redis.EXPECT().Set(gomock.Any(), token, name, gomock.Any()).Return(nil).Times(redisRun)

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
		times int
		err   error
	}
	reduisRun struct {
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
			times int
			err   error
		}{
			times: 0,
			err:   nil,
		},
		reduisRun: struct {
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
			times int
			err   error
		}{
			times: 1,
			err:   errors.New("auth error"),
		},
		reduisRun: struct {
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
			times int
			err   error
		}{
			times: 1,
			err:   nil,
		},
		reduisRun: struct {
			times int
			err   error
		}{
			times: 1,
			err:   errors.New("redis error"),
		},
	},
}

func TestAuthService_Login_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	userService := mock.NewMockuserService(ctrl)
	auth := mock.NewMockauthTool(ctrl)
	reduis := mock.NewMockredis(ctrl)

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
				Return(token, c.authRun.err).
				Times(c.authRun.times)
			reduis.
				EXPECT().
				Set(gomock.Any(), token, c.params.Name, gomock.Any()).
				Return(c.reduisRun.err).
				Times(c.reduisRun.times)

			s := New(userService, reduis, auth)
			got, err := s.Login(context.Background(), c.params)
			if got != nil || err == nil {
				t.Fatal("should get error")
			}
			// shuld be ordering
			if c.reduisRun.err != nil {
				if !errors.Is(err, c.reduisRun.err) {
					t.Errorf("want: %v, got: %v", c.reduisRun.err, err)
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
