package auth

import (
	"context"
	"errors"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

const redis_Error_mes = "something went wrong, redis can not found the token"

var createToken = func(ctx context.Context, name string, auth authTool, redisSvc redisSvc) (*model.Authentication, error) {
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
		ExpiredTime: expired,
	}, nil
}

var verifyToken = func(ctx context.Context, name, token string, auth authTool, redisSvc redisSvc) (bool, error) {
	ret, rName, err := auth.VerifyToken(token)
	if !ret {
		return false, err
	}

	redName, err := redisSvc.Get(ctx, token)
	if err != nil {
		return false, errors.New(redis_Error_mes)
	}

	if redName != rName {
		return false, errors.New("check vaild")
	}

	if name != "" && name != redName {
		return false, errors.New(redis_Error_mes)
	}

	return true, nil
}
