package account

import (
	"context"
	"fmt"
	"net"

	"github.com/231031/ecom-mcs-grpc/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	pb.UnimplementedAccountServiceServer
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(
		serv,
		&grpcServer{
			service: s,
		},
	)
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostAccountSeller(ctx context.Context, r *pb.PostAccountSellerRequest) (*pb.PostAccountSellerResponse, error) {
	seller := Seller{
		StoreName: r.StoreName,
		BaseInfo: BaseInfo{
			FirstName: r.BaseInfo.FirstName,
			LastName:  r.BaseInfo.LastName,
			Email:     r.BaseInfo.Email,
			Phone:     r.BaseInfo.Phone,
			Address:   r.BaseInfo.Address,
		},
	}
	a, err := s.service.PostAccountSeller(ctx, seller)
	if err != nil {
		return nil, err
	}

	return &pb.PostAccountSellerResponse{Account: &pb.AccountSeller{
		Id:        a.ID,
		StoreName: a.StoreName,
		BaseInfo: &pb.BaseInfo{
			FirstName: a.BaseInfo.FirstName,
			LastName:  a.BaseInfo.LastName,
			Email:     a.BaseInfo.Email,
			Phone:     a.BaseInfo.Phone,
			Address:   a.BaseInfo.Address,
		},
	}}, nil
}

func (s *grpcServer) PostAccountBuyer(ctx context.Context, r *pb.PostAccountBuyerRequest) (*pb.PostAccountBuyerResponse, error) {
	buyer := Buyer{
		BaseInfo: BaseInfo{
			FirstName: r.BaseInfo.FirstName,
			LastName:  r.BaseInfo.LastName,
			Email:     r.BaseInfo.Email,
			Password:  r.Password,
			Phone:     r.BaseInfo.Phone,
			Address:   r.BaseInfo.Address,
		},
	}
	a, err := s.service.PostAccountBuyer(ctx, buyer)
	if err != nil {
		return nil, err
	}

	return &pb.PostAccountBuyerResponse{Account: &pb.AccountBuyer{
		Id: a.ID,
		BaseInfo: &pb.BaseInfo{
			FirstName: a.BaseInfo.FirstName,
			LastName:  a.BaseInfo.LastName,
			Email:     a.BaseInfo.Email,
			Phone:     a.BaseInfo.Phone,
			Address:   a.BaseInfo.Address,
		},
	}}, nil
}

func (s *grpcServer) UpdateAccountBuyer(ctx context.Context, a *pb.AccountBuyer) (*pb.AccountBuyer, error) {
	buyer := Buyer{
		ID: a.Id,
		BaseInfo: BaseInfo{
			FirstName: a.BaseInfo.FirstName,
			LastName:  a.BaseInfo.LastName,
			Phone:     a.BaseInfo.Phone,
			Address:   a.BaseInfo.Address,
		},
	}

	_, err := s.service.UpdateAccountBuyer(ctx, buyer)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *grpcServer) UpdateAccountSeller(ctx context.Context, a *pb.AccountSeller) (*pb.AccountSeller, error) {
	seller := Seller{
		ID: a.Id,
		BaseInfo: BaseInfo{
			FirstName: a.BaseInfo.FirstName,
			LastName:  a.BaseInfo.LastName,
			Phone:     a.BaseInfo.Phone,
			Address:   a.BaseInfo.Address,
		},
		StoreName: a.StoreName,
	}

	_, err := s.service.UpdateAccountSeller(ctx, seller)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountBuyerRequest) (*pb.GetAccountBuyerResponse, error) {
	a, err := s.service.GetAccountBuyer(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountBuyerResponse{Account: &pb.AccountBuyer{
		Id: a.ID,
	}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsSellerRequest) (*pb.GetAccountsSellerResponse, error) {
	a, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}

	accounts := []*pb.AccountSeller{}
	for _, p := range a {
		accounts = append(
			accounts,
			&pb.AccountSeller{Id: p.ID},
		)
	}
	return &pb.GetAccountsSellerResponse{Accounts: accounts}, nil
}
