package auth

import (
	"context"
	"errors"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

const redis_Error_mes = "something went wrong, reduis can not found the token"

type redis interface {
	Get(ctx context.Context, token string) (string, error)
	Set(ctx context.Context, key string, value string, expiration int64) error
	Delete(ctx context.Context, key string) error
}

type userService interface {
	Authenticate(ctx context.Context, username string, password string) (bool, error)
}

type authTool interface {
	CreateToken(name string) (string, int64, error)
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
	// TODO with retry
	success, err := s.userServer.Authenticate(ctx, params.Name, params.Password)
	if err != nil {
		return nil, err
	}

	if !success {
		return nil, nil
	}

	return CreateToken(ctx, params.Name, s.auth, s.redis)
}

func (s *AuthService) Verify(ctx context.Context, params *model.Authentication) (bool, error) {
	return Verify(ctx, params.Name, params.Token, s.auth, s.redis)
}

func (s *AuthService) Refresh(ctx context.Context, params *model.Authentication) (*model.Authentication, error) {
	v, err := Verify(ctx, params.Name, params.Token, s.auth, s.redis)
	if err != nil {
		return nil, err
	}

	if !v {
		return nil, errors.New("not vaild token")
	}

	return CreateToken(ctx, params.Name, s.auth, s.redis)
}
