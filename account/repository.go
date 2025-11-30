package account

import (
	"context"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("account not found")
)

type Repository interface {
	Close() error
	CreateAccountSeller(ctx context.Context, a Seller) error
	CreateAccountBuyer(ctx context.Context, a Buyer) error

	UpdateAccountSeller(ctx context.Context, a Seller) error
	UpdateAccountBuyer(ctx context.Context, a Buyer) error

	GetSellerByID(ctx context.Context, id string) (*Seller, error)
	GetBuyerByID(ctx context.Context, id string) (*Buyer, error)

	ListAccountSellers(ctx context.Context, skip uint64, take uint64) ([]Seller, error)
	ListAccountSellersByID(ctx context.Context, ids []string) ([]Seller, error)
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(url string) (Repository, error) {

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() error {
	postgres, err := r.db.DB()
	if err != nil {
		return err
	}
	return postgres.Close()
}

func (r *postgresRepository) Ping() error {
	postgres, err := r.db.DB()
	if err != nil {
		return err
	}
	return postgres.Ping()
}

func (r *postgresRepository) CreateAccountSeller(ctx context.Context, a Seller) error {
	result := r.db.WithContext(ctx).Create(&a)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *postgresRepository) CreateAccountBuyer(ctx context.Context, a Buyer) error {
	result := r.db.WithContext(ctx).Create(&a)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *postgresRepository) UpdateAccountSeller(ctx context.Context, a Seller) error {
	result := r.db.WithContext(ctx).Model(&Seller{}).Where("id = ?", a.ID).Updates(a)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *postgresRepository) UpdateAccountBuyer(ctx context.Context, a Buyer) error {
	result := r.db.WithContext(ctx).Model(&Buyer{}).Where("id = ?", a.ID).Updates(a)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *postgresRepository) GetSellerByID(ctx context.Context, id string) (*Seller, error) {

	var seller = Seller{ID: id}
	result := r.db.WithContext(ctx).First(&seller)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}

	return &seller, nil
}

func (r *postgresRepository) GetBuyerByID(ctx context.Context, id string) (*Buyer, error) {
	var buyer = Buyer{ID: id}
	result := r.db.WithContext(ctx).First(&buyer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}

	return &buyer, nil
}

func (r *postgresRepository) ListAccountSellers(ctx context.Context, skip uint64, take uint64) ([]Seller, error) {
	var sellers = []Seller{}
	result := r.db.WithContext(ctx).Limit(int(take)).Offset(int(skip)).Find(&sellers)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}

	return sellers, nil
}

func (r *postgresRepository) ListAccountSellersByID(ctx context.Context, ids []string) ([]Seller, error) {
	var sellers = []Seller{}
	result := r.db.WithContext(ctx).Where("(id) IN (?)", ids).Find(&sellers)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}

	return sellers, nil
}
