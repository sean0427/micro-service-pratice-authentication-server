package auth

import (
	"context"
	"errors"
	"time"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

type redisSvc interface {
	Get(ctx context.Context, token string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Time) error
	Delete(ctx context.Context, key string) error
}

type userService interface {
	Authenticate(ctx context.Context, username string, password string) (bool, error)
}

type authTool interface {
	CreateToken(name string) (string, time.Time, error)
	VerifyToken(token string) (bool, string, error)
}

type AuthService struct {
	redis      redisSvc
	userServer userService
	auth       authTool
}

func New(user userService, redisRepo redisSvc, auth authTool) *AuthService {
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

	return createToken(ctx, params.Name, s.auth, s.redis)
}

func (s *AuthService) Verify(ctx context.Context, params *model.Authentication) (bool, error) {
	return verifyToken(ctx, params.Name, params.Token, s.auth, s.redis)
}

func (s *AuthService) Refresh(ctx context.Context, params *model.Authentication) (*model.Authentication, error) {
	v, err := verifyToken(ctx, params.Name, params.Token, s.auth, s.redis)
	if err != nil {
		return nil, err
	}

	if !v {
		return nil, errors.New("not vaild token")
	}

	return createToken(ctx, params.Name, s.auth, s.redis)
}
