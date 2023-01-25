package redis_repo

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/sean0427/micro-service-pratice-auth-domain/mock"
)

func FuzzRedisRepo_Get(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add("key"+strconv.Itoa(i), strconv.Itoa(i)+"rest", "")
	}
	for i := 0; i < 100; i++ {
		f.Add("any", "any_retun1", "err"+strconv.Itoa(i))
	}

	f.Fuzz(func(t *testing.T, key, returnedKey, retrunErrMsg string) {
		ctrl := gomock.NewController(t)
		redisClient := mock.NewMockredisClient(ctrl)

		redisClient.EXPECT().
			Get(gomock.Any(), key).
			DoAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
				var err error = nil
				if retrunErrMsg != "" {
					err = errors.New(retrunErrMsg)
				}
				res := redis.NewStringCmd(ctx)

				res.SetErr(err)
				res.SetVal(returnedKey)

				return res
			}).
			Times(1)

		r := New(redisClient)

		got, err := r.Get(context.Background(), key)
		if retrunErrMsg != "" {
			if err == nil {
				t.Errorf("Expect error = %v, but nil", err)
			}
			return
		}
		if got != returnedKey {
			t.Errorf("expect %s but %s", returnedKey, got)
		}
	})
}

func FuzzRedisRepo_Set(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add("key"+strconv.Itoa(i), strconv.Itoa(i)+"val", int64(i*100), true, "")
	}
	for i := 0; i < 100; i++ {
		f.Add("any1", strconv.Itoa(i)+"va1l", int64(i*100), false, "")
	}
	for i := 0; i < 100; i++ {
		f.Add("any", strconv.Itoa(i)+"val", int64(i*100), false, "err")
	}

	f.Fuzz(func(t *testing.T, key, value string, exp int64, returnedSuccess bool, errMsg string) {
		ctrl := gomock.NewController(t)
		redisClient := mock.NewMockredisClient(ctrl)
		redisClient.EXPECT().
			SetNX(gomock.Any(), key, value, gomock.Any()).
			DoAndReturn(func(ctx context.Context, key string, value string, expiration time.Duration) *redis.BoolCmd {
				var err error = nil
				if errMsg != "" {
					err = errors.New(errMsg)
				}

				res := redis.NewBoolCmd(ctx)
				res.SetVal(returnedSuccess)
				res.SetErr(err)

				return res
			}).
			Times(1)

		r := New(redisClient)
		err := r.Set(context.Background(), key, value, exp)

		if errMsg != "" {
			if err == nil {
				t.Errorf("Expect error = %v, but nil", err)
			}
			return
		}
		if !returnedSuccess {
			if err == nil {
				t.Errorf("Expect error = %v, but nil", err)
			}
			return
		}

		if err != nil {
			t.Errorf("Expect error nil, but %v", err)
		}
	})
}

func FuzzRedisRepo_Delete(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add("key"+strconv.Itoa(i), "")
	}
	for i := 0; i < 100; i++ {
		f.Add("any", "err"+strconv.Itoa(i))
	}

	f.Fuzz(func(t *testing.T, key, retrunErrMsg string) {
		ctrl := gomock.NewController(t)
		redisClient := mock.NewMockredisClient(ctrl)

		redisClient.EXPECT().
			Del(gomock.Any(), key).
			DoAndReturn(func(ctx context.Context, key string) *redis.IntCmd {
				var err error = nil
				if retrunErrMsg != "" {
					err = errors.New(retrunErrMsg)
				}
				res := redis.NewIntCmd(ctx)

				res.SetErr(err)

				return res
			}).
			Times(1)

		r := New(redisClient)

		err := r.Delete(context.Background(), key)
		if retrunErrMsg != "" {
			if err == nil {
				t.Errorf("Expect error = %v, but nil", err)
			}
			return
		}
	})
}
