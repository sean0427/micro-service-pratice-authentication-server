package auth

import (
	"context"
	"time"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

type redis interface {
	Get(ctx context.Context, token string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Time) error
	Delete(ctx context.Context, key string) error
}

type userService interface {
	Authenticate(ctx context.Context, username string, password string) (bool, error)
}

type authTool interface {
	CreateToken(name string) (string, error)
	VerifyToken(token string) (bool, string)
}

type AuthService struct {
	redis      redis
	userServer userService
	auth       authTool
}

func New(user userService, redisRepo redis, auth authTool) *AuthService {
	return &AuthService{
		redis:      redisRepo,
		userServer: user,
		auth:       auth,
	}
}

func (s *AuthService) Login(ctx context.Context, params *model.LoginInfo) (*model.Authentication, error) {
	success, err := s.userServer.Authenticate(ctx, params.Name, params.Password)
	if err != nil {
		return nil, err
	}

	if !success {
		return nil, nil
	}

	token, err := s.auth.CreateToken(params.Name)
	if err != nil {
		return nil, err
	}

	// TODO
	err = s.redis.Set(ctx, token, params.Name, time.Now().Add(time.Hour*8))
	if err != nil {
		return nil, err
	}

	return &model.Authentication{
		Name:  params.Name,
		Token: token,
	}, nil
}
