package catalog

import (
	"context"
	"log"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64, quantity uint32) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
	UpdateQuantity(ctx context.Context, ids []string, quantity []uint32) ([]string, error)
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{repository: r}
}
func (s *catalogService) PostProduct(ctx context.Context, name, description string, price float64, quantity uint32) (*Product, error) {
	p := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().String(),
		Quantity:    quantity,
	}

	if err := s.repository.PutProduct(ctx, *p); err != nil {
		log.Println("service: error putting product")

		return nil, err
	}

	return p, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductByID(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListProducts(ctx, skip, take)
}

func (s *catalogService) GetProductByIDs(ctx context.Context, ids []string) ([]Product, error) {
	return s.repository.ListProductsWithIDs(ctx, ids)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProducts(ctx, query, skip, take)
}

func (s *catalogService) UpdateQuantity(ctx context.Context, ids []string, quantity []uint32) ([]string, error) {
	err := s.repository.UpdateQuantity(ctx, ids, quantity)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
