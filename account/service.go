package account

import (
	"context"

	"github.com/231031/ecom-mcs-grpc/account/utils"
	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccountSeller(ctx context.Context, a Seller) (*Seller, error)
	PostAccountBuyer(ctx context.Context, a Buyer) (*Buyer, error)
	UpdateAccountSeller(ctx context.Context, a Seller) (*Seller, error)
	UpdateAccountBuyer(ctx context.Context, a Buyer) (*Buyer, error)
	GetAccountBuyer(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AccountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &AccountService{repository: r}
}

func (s *AccountService) PostAccountSeller(ctx context.Context, a Seller) (*Seller, error) {
	a.ID = ksuid.New().String()

	hashed, err := utils.HashPassword(a.Password)
	if err != nil {
		return nil, err
	}

	a.Password = hashed
	if err := s.repository.CreateAccountSeller(ctx, a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AccountService) PostAccountBuyer(ctx context.Context, a Buyer) (*Buyer, error) {
	a.ID = ksuid.New().String()

	hashed, err := utils.HashPassword(a.Password)
	if err != nil {
		return nil, err
	}

	a.Password = hashed
	if err := s.repository.CreateAccountBuyer(ctx, a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AccountService) UpdateAccountBuyer(ctx context.Context, a Buyer) (*Buyer, error) {
	err := s.repository.UpdateAccountBuyer(ctx, a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *AccountService) UpdateAccountSeller(ctx context.Context, a Seller) (*Seller, error) {
	err := s.repository.UpdateAccountSeller(ctx, a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *AccountService) GetAccountBuyer(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *AccountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}
