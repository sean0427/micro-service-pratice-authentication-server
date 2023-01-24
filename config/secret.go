package config

import (
	"os"
)

const defaultJWTKey = ".secret/jwt.key"

func GetJWTSecretKey() string {
	path, found := os.LookupEnv("JWT_SECRET_KEY_FILE")
	if !found {
		path = defaultJWTKey
	}
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b)
}
