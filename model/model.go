package model

import (
	"time"
)

type Authentication struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Token       string    `json:"token"`
	ExpiredTime time.Time `json:"expired_time"`
	Timestamp   time.Time `json:"timestamp"`
}

type LoginInfo struct {
	Password string `json:"password"`
	Name     string `json:"name"`
}
