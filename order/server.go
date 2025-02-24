package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/231031/ecom-mcs-grpc/account"
	"github.com/231031/ecom-mcs-grpc/catalog"
	"github.com/231031/ecom-mcs-grpc/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ErrInvalidOrder   = errors.New("failed to create order")
	ErrInvalidAccount = errors.New("account not found")
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serve := grpc.NewServer()
	grpcServiceServer := &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	}

	pb.RegisterOrderServiceServer(serve, grpcServiceServer)
	reflection.Register(serve)
	return serve.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("error getting account", err)
		return nil, ErrInvalidAccount
	}

	productIDs := []string{}
	for _, p := range r.Products {
		if p.Quantity != 0 {
			productIDs = append(productIDs, p.ProductId)
		}
	}

	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("error getting product", err)
		return nil, err
	}

	if len(products) != len(productIDs) {
		notFound := len(productIDs) - len(products)
		return nil, fmt.Errorf("%d products not found", notFound)
	}

	ids := []string{}
	quantity := []uint32{}
	orderProducts := []OrderedProduct{}
	for _, p := range products {
		product := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}
		for _, rp := range r.Products {
			if product.ID == rp.ProductId {
				product.Quantity = rp.Quantity

				ids = append(ids, p.ID)
				quantity = append(quantity, p.Quantity-rp.Quantity)
				break
			}
		}
		orderProducts = append(orderProducts, product)
	}

	order, err := s.service.PostOrder(ctx, r.AccountId, orderProducts)
	if err != nil {
		log.Println(err)
		return nil, ErrInvalidOrder
	}

	// update quantity of products
	_, err = s.catalogClient.UpdateQuantity(ctx, ids, quantity)
	if err != nil {
		return nil, err
	}

	productsPsroto := []*pb.Order_OrderProduct{}
	for _, p := range order.Products {
		productsPsroto = append(productsPsroto, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	createdAtBin, err := order.CreatedAt.MarshalBinary()
	if err != nil {
		log.Println("error marshalling time", err)
		return nil, errors.New("could not marshal timestamp")
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		CreatedAt:  createdAtBin,
		Products:   productsPsroto,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
	}

	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrderForAccountRequest) (*pb.GetOrderForAccountResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("error getting account", err)
		return nil, ErrInvalidAccount
	}

	orders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("error getting orders", err)
		return nil, err
	}

	productIDsMap := map[string]bool{}
	for _, o := range orders {
		for _, p := range o.Products {
			productIDsMap[p.ID] = true
		}
	}

	productIDs := []string{}
	for id := range productIDsMap {
		productIDs = append(productIDs, id)
	}
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("error getting products", err)
		return nil, err
	}

	ordersProto := []*pb.Order{}
	for _, o := range orders {
		order := &pb.Order{
			Id:         o.ID,
			AccountId:  o.AccountID,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}
		order.CreatedAt, err = o.CreatedAt.MarshalBinary()
		if err != nil {
			log.Println("error marshal timestamp", err)
		}

		for _, p := range o.Products {
			for _, queryP := range products {
				if queryP.ID == p.ID {
					product := &pb.Order_OrderProduct{
						Id:          queryP.ID,
						Name:        queryP.Name,
						Description: queryP.Description,
						Price:       queryP.Price,
						Quantity:    p.Quantity,
					}
					order.Products = append(order.Products, product)
					break
				}
			}
		}

		ordersProto = append(ordersProto, order)
	}

	return &pb.GetOrderForAccountResponse{
		Orders: ordersProto,
	}, nil

}
