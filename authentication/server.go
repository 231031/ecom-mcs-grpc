package authentication

import (
	"context"
	"fmt"
	"net"

	"github.com/231031/ecom-mcs-grpc/authentication/model"
	"github.com/231031/ecom-mcs-grpc/authentication/pb"
	"github.com/231031/ecom-mcs-grpc/authentication/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service service.Service
	pb.UnimplementedAuthenticationServiceServer
}

func ListenGRPC(s service.Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAuthenticationServiceServer(
		serv,
		&grpcServer{
			service: s,
		},
	)
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	u := &model.User{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
	_, err := s.service.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{}, nil
}

func (s *grpcServer) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.TokenResponse, error) {
	token, err := s.service.LoginUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, err
	}

	tokenResp := &pb.TokenResponse{
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return tokenResp, nil
}

func (s *grpcServer) RefreshTokenUser(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	token, err := s.service.RefreshTokenUser(ctx, in.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	tokenResp := &pb.TokenResponse{
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return tokenResp, nil
}
