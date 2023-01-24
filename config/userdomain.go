package config

import (
	"errors"
	"os"
)

func GetUserAuthGrpcAddr() (string, error) {
	if v, found := os.LookupEnv(""); found {
		return v, nil
	} else {
		return "", errors.New("USER_AUTHGRPC_ADDR is not set")
	}
}
