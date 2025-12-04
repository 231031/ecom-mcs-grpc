package authentication

import (
	"context"

	"github.com/231031/ecom-mcs-grpc/authentication/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AuthenticationServiceClient
}

func NewClient(url string, opts ...grpc.DialOption) (*Client, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	finalOpts := append(defaultOpts, opts...)

	conn, err := grpc.NewClient(
		url,
		finalOpts...,
	)
	if err != nil {
		return nil, err
	}

	c := pb.NewAuthenticationServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	u, err := c.service.CreateUser(ctx, in)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *Client) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	u, err := c.service.LoginUser(ctx, in)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *Client) RefreshTokenUser(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	token, err := c.service.RefreshTokenUser(ctx, in)
	if err != nil {
		return nil, err
	}
	return token, nil
}
