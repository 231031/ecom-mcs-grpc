package authentiction

import (
	"fmt"
	"net"

	"github.com/231031/ecom-mcs-grpc/authentiction/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	pb.UnimplementedAuthenticationServiceServer
}

func ListenGRPC(s Service, port int) error {
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
