package catalog

import (
	"context"
	"fmt"
	"net"

	"github.com/231031/ecom-mcs-grpc/account"
	"github.com/231031/ecom-mcs-grpc/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	pb.UnimplementedCatalogServiceServer
}

func ListenGRPC(s Service, port int, accountURL string) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		return err
	}

	serve := grpc.NewServer()
	grpcServiceServer := &grpcServer{
		service:       s,
		accountClient: accountClient,
	}

	// register server with pb
	pb.RegisterCatalogServiceServer(serve, grpcServiceServer)
	reflection.Register(serve)
	return serve.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	_, err := s.accountClient.GetAccountSellerByID(ctx, r.SellerId)
	if err != nil {
		return nil, err
	}

	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.SellerId, r.Price, r.Quantity)
	if err != nil {
		return nil, err
	}

	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			SellerId:    p.SellerID,
		},
	}, nil

}
func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			SellerId:    p.SellerID,
		},
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res []Product
	var err error

	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) > 0 {
		res, err = s.service.GetProductByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}

	if err != nil {
		return nil, err
	}

	products := []*pb.Product{}
	for _, p := range res {
		products = append(products, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}
	return &pb.GetProductsResponse{Products: products}, nil
}

func (s *grpcServer) UpdateQuantity(ctx context.Context, req *pb.UpdateQuantityRequest) (*pb.UpdateQuantityResponse, error) {
	ids, err := s.service.UpdateQuantity(ctx, req.Ids, req.Quantity)
	if err != nil {
		return nil, err
	}

	idsResp := &pb.UpdateQuantityResponse{
		Ids: ids,
	}
	return idsResp, nil
}

func (s *grpcServer) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	p := Product{
		ID:          req.Id,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Quantity:    req.Quantity,
	}

	_, err := s.service.UpdateProduct(ctx, p)
	if err != nil {
		return nil, err
	}

	return req, nil
}
