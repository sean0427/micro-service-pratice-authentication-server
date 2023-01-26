package auth

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/sean0427/micro-service-pratice-auth-domain/mock"
	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

func TestMain(m *testing.M) {
	orgVerify := verifyToken
	code := m.Run()
	verifyToken = orgVerify
	os.Exit(code)
}

func FuzzAuthService_Login(f *testing.F) {
	orgCreateToken := createToken

	f.Add("name_any", "pw", "token_any", true)
	f.Add("deaff", "fea", "awefafwefeawf", true)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", true)
	f.Add("deaff", "fea", "awefafwefeawf", false)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", false)
	f.Add("123123132", "3123", "123131fewafeawfeaffw", false)

	f.Fuzz(func(t *testing.T, name, pw, token string, success bool) {
		ctrl := gomock.NewController(t)
		userService := mock.NewMockuserService(ctrl)

		userService.EXPECT().Authenticate(gomock.Any(), name, pw).Return(success, nil).Times(1)
		createToken = func(ctx context.Context, inName string, _ authTool, _ redisSvc) (*model.Authentication, error) {
			return &model.Authentication{
				Name:  inName,
				Token: token,
			}, nil
		}

		s := New(userService, nil, nil)
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

	createToken = orgCreateToken
}

var testProductService_Login_Error = []struct {
	name             string
	params           *model.LoginInfo
	userSericeReturn struct {
		success bool
		err     error
	}
	createTokenRun struct {
		err     error
		expired int64
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
		createTokenRun: struct {
			err     error
			expired int64
		}{
			err: nil,
		},
	},
	{
		name: "createToken return error",
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
		createTokenRun: struct {
			err     error
			expired int64
		}{
			err: errors.New("auth error"),
		},
	},
}

func TestAuthService_Login_Error(t *testing.T) {
	orgCreateToken := createToken

	ctrl := gomock.NewController(t)
	userService := mock.NewMockuserService(ctrl)
	auth := mock.NewMockauthTool(ctrl)
	redis := mock.NewMockredisSvc(ctrl)

	for _, c := range testProductService_Login_Error {
		t.Run(c.name, func(t *testing.T) {
			createToken = func(_ context.Context, inName string, _ authTool, _ redisSvc) (*model.Authentication, error) {
				return nil, c.createTokenRun.err
			}

			userService.
				EXPECT().
				Authenticate(gomock.Any(), c.params.Name, c.params.Password).
				Return(c.userSericeReturn.success, c.userSericeReturn.err).
				Times(1)

			s := New(userService, redis, auth)
			got, err := s.Login(context.Background(), c.params)
			if got != nil || err == nil {
				t.Fatal("should get error")
			}

			if c.createTokenRun.err != nil {
				if !errors.Is(err, c.createTokenRun.err) {
					t.Errorf("want: %v, got: %v", c.createTokenRun.err, err)
				}
				return
			}

			if !errors.Is(err, c.userSericeReturn.err) {
				t.Errorf("want: %v, got: %v", c.userSericeReturn.err, err)
			}
		})
	}

	createToken = orgCreateToken
}

func TestAuthService_Verify(t *testing.T) {
	testParams := &model.Authentication{}
	testError := errors.New("any error")

	orgVerify := verifyToken
	t.Run("test verify", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		auth := mock.NewMockauthTool(ctrl)
		redis := mock.NewMockredisSvc(ctrl)

		s := New(nil, redis, auth)

		verifyToken = func(ctx context.Context, name, token string, _ authTool, _ redisSvc) (bool, error) {
			if name != testParams.Name {
				t.Errorf("want: %v, got: %v", testParams.Name, name)
			}

			if token != testParams.Token {
				t.Errorf("want: %v, got: %v", testParams.Token, token)
			}

			return true, testError
		}

		got, err := s.Verify(context.Background(), testParams)

		if !errors.Is(testError, err) {
			t.Errorf("want: %v, got: %v", testError, err)
		}

		if got != true {
			t.Errorf("AuthService.Verify() = %v, want %v", got, true)
		}
	})

	verifyToken = orgVerify
}

func FuzzAuthService_Refresh(f *testing.F) {
	orgVerify := verifyToken
	orgCreateToken := createToken

	f.Add("test", "fjieajfioef", "efaefef", true, "", "")
	f.Add("test2", "fejifeafjaeif", "eafeafewf", false, "", "")
	f.Add("test3", "fdefejiafjaeif", "fawefee124e1", false, "feaffewaf", "")
	f.Add("test3", "fejiafjeaeif", "feawfef", false, "", "faewfewaf")

	f.Fuzz(func(t *testing.T, name, token, newToken string, verifyTokeReturn bool, verifyErrMsg, createErroMsg string) {
		s := &AuthService{}

		verifyRunTimes := 0
		createTokenTimes := 0
		verifyToken = func(ctx context.Context, inName, inToken string, auth authTool, redisSvc redisSvc) (bool, error) {
			verifyRunTimes++
			if inName != name {
				t.Errorf("want: %v, got: %v", name, inName)
			}

			if inToken != token {
				t.Errorf("want: %v, got: %v", token, inToken)
			}

			if verifyErrMsg != "" {
				return false, errors.New(verifyErrMsg)
			}

			return verifyTokeReturn, nil
		}

		createToken = func(_ context.Context, inName string, auth authTool, redisSvc redisSvc) (*model.Authentication, error) {
			createTokenTimes++

			if inName != name {
				t.Errorf("want: %v, got: %v", name, inName)
			}

			if createErroMsg != "" {
				return nil, errors.New(verifyErrMsg)
			}

			return &model.Authentication{
				Name:  name,
				Token: newToken,
			}, nil
		}

		got, err := s.Refresh(context.Background(), &model.Authentication{
			Name:  name,
			Token: token,
		})
		if verifyRunTimes != 1 {
			t.Errorf("verify should be always run once")
		}

		if !verifyTokeReturn {
			if err == nil {
				t.Error("expected error")
			}

			if createTokenTimes > 0 {
				t.Error("create token should not have been called")
			}
			return
		}

		if createErroMsg != "" || verifyErrMsg != "" {
			if err == nil {
				t.Error("expected error")
			}
			return
		}

		if got.Token != newToken {
			t.Errorf("want: %v, got: %v", newToken, got.Token)
		}
	})
	verifyToken = orgVerify
	createToken = orgCreateToken

}
