package redis_repo

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type RedisRepo struct {
	// client *redis.Client
	temp map[string]string // using string for first implement
}

func New(client *redis.Client) *RedisRepo {
	return &RedisRepo{
		// client: client,
	}
}

func (r *RedisRepo) Get(ctx context.Context, key string) (string, error) {
	v := r.temp[key]
	return v, nil
}

func (r *RedisRepo) Set(ctx context.Context, key, value string, expireation int64) error {
	r.temp[key] = value
	// TODO: expire
	return nil
}

func (r *RedisRepo) Delete(ctx context.Context, key string) error {
	delete(r.temp, key)
	return nil
}
