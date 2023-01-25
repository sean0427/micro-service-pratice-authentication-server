package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v9"
	service "github.com/sean0427/micro-service-pratice-auth-domain"
	"github.com/sean0427/micro-service-pratice-auth-domain/config"
	jwt_token_helper "github.com/sean0427/micro-service-pratice-auth-domain/jwttokenhelper"
	handler "github.com/sean0427/micro-service-pratice-auth-domain/net"
	"github.com/sean0427/micro-service-pratice-auth-domain/redis_repo"
	"github.com/sean0427/micro-service-pratice-auth-domain/userdomainclient"
	pb "github.com/sean0427/micro-service-pratice-auth-domain/userdomainclient/grpc/auth"
	"google.golang.org/grpc"
)

var (
	token_minute = flag.Int("token-period", 36, "access-token minute preiod")
	port         = flag.Int("port", 8080, "port")
)

func createGrpcClient(addr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getRedisClient() (*redis.Client, error) {
	address, err := config.GetRedisAddress()
	if err != nil {
		return nil, err
	}

	password, err := config.GetRedisPassword()
	if err != nil {
		return nil, err
	}

	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	}), nil
}

func startServer() {
	fmt.Println("Starting server...")

	addr, err := config.GetUserAuthGrpcAddr()
	if err != nil {
		panic(err)
	}

	conn, err := createGrpcClient(addr)
	if err != nil {
		panic(err)
	}

	userClient := pb.NewAuthClient(conn)
	defer conn.Close()
	userDomainClient := userdomainclient.New(userClient)

	authHelper := jwt_token_helper.New([]byte(config.GetJWTSecretKey()),
		time.Minute*time.Duration(*token_minute))

	rdb, err := getRedisClient()
	if err != nil {
		panic(err)
	}
	redis := redis_repo.New(rdb)

	s := service.New(userDomainClient, redis, authHelper)
	h := handler.New(s)

	handler := h.Handler()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), handler); err != nil {
		panic(err)
	}

	fmt.Println("Stoping server...")
}

func main() {
	startServer()
}
