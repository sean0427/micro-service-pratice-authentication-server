package auth

import (
	"context"
	"errors"
	"time"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

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

var Verify = func(ctx context.Context, name, token string, auth authTool, redisSvc redis) (bool, error) {
	if v, msg := auth.VerifyToken(token); !v {
		return false, errors.New(msg)
	}

	v, err := redisSvc.Get(ctx, token)
	if err != nil || v != name {
		return false, errors.New(redis_Error_mes)
	}

	return true, nil
}
