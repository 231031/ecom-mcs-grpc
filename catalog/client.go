package catalog

import (
	"context"
	"errors"

	"github.com/231031/ecom-mcs-grpc/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrNotHaveProductsInfo = errors.New("product info is not available")
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewCatalogServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description, seller_id string, price float64, quantity uint32) (*Product, error) {
	r, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		SellerId:    seller_id,
	})
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
		Quantity:    r.Product.Quantity,
		SellerID:    r.Product.SellerId,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
		Quantity:    r.Product.Quantity,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]Product, error) {
	r, err := c.service.GetProducts(
		ctx,
		&pb.GetProductsRequest{
			Skip:  skip,
			Take:  take,
			Query: query,
			Ids:   ids,
		},
	)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for _, p := range r.Products {
		products = append(products, Product{
			ID:          p.Id,
			Description: p.Description,
			Name:        p.Name,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	return products, nil
}

func (c *Client) UpdateQuantity(ctx context.Context, ids []string, quantity []uint32) ([]string, error) {
	if len(ids) == 0 || len(quantity) == 0 {
		return nil, ErrNotHaveProductsInfo
	}

	resp, err := c.service.UpdateQuantity(
		ctx,
		&pb.UpdateQuantityRequest{
			Ids:      ids,
			Quantity: quantity,
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.Ids, nil
}

func (c *Client) UpdateProduct(ctx context.Context, id, name, description string, price float64, quantity uint32) (*Product, error) {
	p, err := c.service.UpdateProduct(ctx, &pb.Product{
		Id:          id,
		Name:        name,
		Price:       price,
		Description: description,
		Quantity:    quantity,
	})
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
	}, nil
}
