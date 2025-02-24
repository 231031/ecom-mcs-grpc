package order

import (
	"context"
	"log"
	"time"

	"github.com/231031/ecom-mcs-grpc/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	service := pb.NewOrderServiceClient(conn)

	return &Client{
		conn:    conn,
		service: service,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	productsProto := []*pb.PostOrderRequest_OrderProduct{}
	for _, p := range products {
		productsProto = append(productsProto, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}

	orderProto, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountID,
		Products:  productsProto,
	})
	if err != nil {
		return nil, err
	}

	order := &Order{
		ID:         orderProto.Order.Id,
		TotalPrice: orderProto.Order.TotalPrice,
		AccountID:  orderProto.Order.AccountId,
		Products:   products,
	}

	createdAt := time.Time{}
	err = createdAt.UnmarshalBinary(orderProto.Order.CreatedAt)
	if err != nil {
		log.Println("error unmarshalling timestamp", err)
	}
	order.CreatedAt = createdAt

	return order, nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	req := &pb.GetOrderForAccountRequest{
		AccountId: accountID,
	}
	ordersProto, err := c.service.GetOrdersForAccount(ctx, req)
	if err != nil {
		log.Println("error getting orders", err)
		return nil, err
	}

	orders := []Order{}
	for _, op := range ordersProto.Orders {
		createdAt := time.Time{}
		err = createdAt.UnmarshalBinary(op.CreatedAt)
		if err != nil {
			log.Println("error unmarshalling timestamp", err)
		}

		order := Order{
			ID:         op.Id,
			AccountID:  op.AccountId,
			TotalPrice: op.TotalPrice,
			CreatedAt:  createdAt,
		}

		products := []OrderedProduct{}
		for _, p := range op.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}

		order.Products = products
		orders = append(orders, order)
	}

	return orders, nil
}
