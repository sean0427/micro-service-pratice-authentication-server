package userdomainclient

import (
	"context"
	"errors"

	pb "github.com/sean0427/micro-service-pratice-auth-domain/userdomainclient/grpc/auth"
	"google.golang.org/grpc"
)

type authClient interface {
	Authenticate(ctx context.Context, in *pb.AuthRequest, opts ...grpc.CallOption) (*pb.AuthReply, error)
}

type GrpcClient struct {
	authClient
}

func New(authClient pb.AuthClient) *GrpcClient {
	return &GrpcClient{
		authClient: authClient,
	}
}

func (c *GrpcClient) Authenticate(ctx context.Context, name, password string) (bool, error) {
	reply, err := c.authClient.Authenticate(ctx, &pb.AuthRequest{Name: name, Password: password})
	if err != nil {
		return false, err
	}

	if reply.Success {
		return true, nil
	}

	return false, errors.New(reply.Message)
}
