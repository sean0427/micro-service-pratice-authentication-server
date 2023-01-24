package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	service "github.com/sean0427/micro-service-pratice-auth-domain"
	"github.com/sean0427/micro-service-pratice-auth-domain/config"
	jwt_token_helper "github.com/sean0427/micro-service-pratice-auth-domain/jwttokenhelper"
	handler "github.com/sean0427/micro-service-pratice-auth-domain/net"
	"github.com/sean0427/micro-service-pratice-auth-domain/redis_repo"
	"github.com/sean0427/micro-service-pratice-auth-domain/userdomainclient"
	pb "github.com/sean0427/micro-service-pratice-auth-domain/userdomainclient/grpc/auth"
	"google.golang.org/grpc"
)

func createGrpcClient(addr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

var (
	token_minute = flag.Int("token-period", 36, "access-token minute preiod")
)

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

	client := pb.NewAuthClient(conn)
	defer conn.Close()

	authHelper := jwt_token_helper.New([]byte(config.GetJWTSecretKey()),
		time.Minute*time.Duration(*token_minute))
	redis := redis_repo.New(nil)
	userDomainClient := userdomainclient.New(client)

	s := service.New(userDomainClient, redis, authHelper)
	h := handler.New(s)

	handler := h.Handler()
	http.ListenAndServe(":8080", handler)

	fmt.Println("Stoping server...")
}

func main() {
	startServer()
}
