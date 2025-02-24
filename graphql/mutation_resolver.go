package graphql

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/231031/ecom-mcs-grpc/order"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *Server
}

func (m *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := m.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Account{ID: a.ID, Name: a.Name}, nil
}

func (m *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := m.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price, uint32(in.Quantity))
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          p.ID,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Quantity:    in.Quantity,
	}, nil
}

func (m *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var products []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, order.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}

	order, err := m.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Order{
		ID:         order.ID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
	}, nil
}
