package redis_repo

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v9"
)

type redisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type RedisRepo struct {
	client redisClient
}

func New(client redisClient) *RedisRepo {
	return &RedisRepo{
		client: client,
	}
}

func (r *RedisRepo) Get(ctx context.Context, key string) (string, error) {
	res := r.client.Get(ctx, key)
	return res.Result()
}

func (r *RedisRepo) Set(ctx context.Context, key, value string, expireation int64) error {
	set := r.client.SetNX(ctx, key, value, time.Until(time.UnixMicro(expireation)))
	if v, err := set.Result(); err != nil {
		return err
	} else if !v {
		return errors.New("Set key failed")
	}
	return nil
}

func (r *RedisRepo) Delete(ctx context.Context, key string) error {
	cmd := r.client.Del(ctx, key)
	_, err := cmd.Result()
	return err
}
