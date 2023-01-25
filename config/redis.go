package config

import (
	"errors"
	"os"
)

func GetRedisAddress() (string, error) {
	if v, found := os.LookupEnv("REDIS_ADDRESS"); found {
		return v, nil
	} else {
		return "", errors.New("REDIS_ADDRESS is not set")
	}
}

func GetRedisPassword() (string, error) {
	if v, found := os.LookupEnv("REDIS_PASSWORD"); found {
		return v, nil
	} else {
		return "", errors.New("REDIS_PASSWORD is not set")
	}
}
