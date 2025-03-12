package graphql

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/231031/ecom-mcs-grpc/account/pb"
	"github.com/231031/ecom-mcs-grpc/order"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrInvalidInfo      = errors.New("invalid info")
)

type mutationResolver struct {
	server *Server
}

func (m *mutationResolver) CreateAccountBuyer(ctx context.Context, in AccountBuyerInput) (*AccountBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	data := &pb.PostAccountBuyerRequest{}
	req := MapGraphQLInputToRequest(in, TypeCreate)
	switch v := req.(type) {
	case *pb.PostAccountBuyerRequest:
		data = v
	default:
		return nil, ErrInvalidInfo
	}

	a, err := m.server.accountClient.PostAccountBuyer(ctx, data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &AccountBuyer{
		ID:        a.ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Email:     a.Email,
		Phone:     a.Phone,
		Address:   a.Address,
		Orders:    []*Order{},
	}, nil
}

func (m *mutationResolver) CreateAccountSeller(ctx context.Context, in AccountSellerInput) (*AccountSeller, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	data := &pb.PostAccountSellerRequest{}
	req := MapGraphQLInputToRequest(in, TypeCreate)
	switch v := req.(type) {
	case *pb.PostAccountSellerRequest:
		data = v
	default:
		return nil, ErrInvalidInfo
	}

	a, err := m.server.accountClient.PostAccountSeller(ctx, data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &AccountSeller{
		ID:        a.ID,
		StoreName: a.StoreName,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Email:     a.Email,
		Phone:     a.Phone,
		Address:   a.Address,
		Products:  []*Product{},
	}, nil
}

func (m *mutationResolver) UpdateAccountBuyer(ctx context.Context, in AccountBuyerInput, id string) (*AccountBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	data := &pb.AccountBuyer{}
	req := MapGraphQLInputToRequest(in, TypeUpdate)
	switch v := req.(type) {
	case *pb.AccountBuyer:
		v.Id = id
		data = v
	default:
		return nil, ErrInvalidInfo
	}

	a, err := m.server.accountClient.UpdateAccountBuyer(ctx, data)
	if err != nil {
		return nil, err
	}

	return &AccountBuyer{
		ID:        a.Id,
		Email:     a.BaseInfo.Email,
		FirstName: a.BaseInfo.FirstName,
		LastName:  a.BaseInfo.LastName,
		Phone:     a.BaseInfo.Phone,
		Address:   a.BaseInfo.Address,
		Orders:    []*Order{},
	}, nil
}

func (m *mutationResolver) UpdateAccountSeller(ctx context.Context, in AccountSellerInput, id string) (*AccountSeller, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	data := &pb.AccountSeller{}
	req := MapGraphQLInputToRequest(in, TypeUpdate)
	switch v := req.(type) {
	case *pb.AccountSeller:
		v.Id = id
		data = v
	default:
		return nil, ErrInvalidInfo
	}

	log.Println("data: ", data)
	a, err := m.server.accountClient.UpdateAccountSeller(ctx, data)
	if err != nil {
		return nil, err
	}

	return &AccountSeller{
		ID:        a.Id,
		Email:     a.BaseInfo.Email,
		FirstName: a.BaseInfo.FirstName,
		LastName:  a.BaseInfo.LastName,
		Phone:     a.BaseInfo.Phone,
		Address:   a.BaseInfo.Address,
		Products:  []*Product{},
	}, nil
}

func (m *mutationResolver) LoginAccount(ctx context.Context, email, password string, role RoleType) (LoginResult, error) {
	if role == RoleTypeSeller {
		return &AccountSeller{
			ID:        "456",
			Email:     email,
			FirstName: "Jane",
			LastName:  "Doe",
			Phone:     "987-654-3210",
			Address:   "456 Elm St",
		}, nil
	}

	// Default: Return an AccountBuyer
	return &AccountBuyer{
		ID:        "123",
		Email:     email,
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "123-456-7890",
		Address:   "123 Main St",
	}, nil
}

func (m *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := m.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.SellerID, in.Price, uint32(in.Quantity))
	if err != nil {
		return nil, err
	}

	log.Println(p.SellerID)
	return &Product{
		ID:          p.ID,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Quantity:    in.Quantity,
		SellerID:    p.SellerID,
	}, nil
}

func (m *mutationResolver) UpdateProduct(ctx context.Context, in ProductInput, id string) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// edit function update product
	_, err := m.server.catalogClient.UpdateProduct(ctx, id, in.Name, in.Description, in.Price, uint32(in.Quantity))
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Quantity:    in.Quantity,
	}, nil
}

func (m *mutationResolver) DeleteProduct(ctx context.Context, id string) (string, error) {
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()
	return id, nil
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
			ID:       p.ProductID,
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

func (m *mutationResolver) DeleteOrder(ctx context.Context, id string) (string, error) {
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()
	return id, nil
}
