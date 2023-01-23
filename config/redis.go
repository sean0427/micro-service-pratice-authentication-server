package config

import (
	"errors"
	"os"
)

func GetRedisPassword() (string, error) {
	if v, found := os.LookupEnv("REDIS_PASSWORD"); found {
		return v, nil
	} else {
		return "", errors.New("REDIS_PASSWORD is not set")
	}
}
