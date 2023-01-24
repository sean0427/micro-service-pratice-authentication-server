package auth

import (
	"context"
	"errors"
	"time"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

const redis_error_mes = "something went wrong, reduis can not found the token"

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

var CreateToken = func(ctx context.Context, name string, auth authTool, redisSvc redis) (*model.Authentication, error) {
	token, expired, err := auth.CreateToken(name)
	if err != nil {
		return nil, err
	}

	err = redisSvc.Set(ctx, token, name, expired)
	if err != nil {
		return nil, err
	}

	return &model.Authentication{
		Name:        name,
		Token:       token,
		ExpiredTime: time.UnixMilli(expired),
	}, nil
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
	if v, msg := s.auth.VerifyToken(params.Token); !v {
		return false, errors.New(msg)
	}

	v, err := s.redis.Get(ctx, params.Token)
	if err != nil || v != params.Name {
		return false, errors.New(redis_error_mes)
	}

	return true, nil
}

func (s *AuthService) Refresh(ctx context.Context, params *model.Authentication) (*model.Authentication, error) {
	v, err := s.Verify(ctx, params)
	if err != nil {
		return nil, err
	}

	if !v {
		return nil, errors.New("not vaild token")
	}

	return CreateToken(ctx, params.Name, s.auth, s.redis)
}
