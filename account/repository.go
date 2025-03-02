package account

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	Close() error
	CreateAccountSeller(ctx context.Context, a Seller) error
	CreateAccountBuyer(ctx context.Context, a Buyer) error

	UpdateAccountSeller(ctx context.Context, a Seller) error
	UpdateAccountBuyer(ctx context.Context, a Buyer) error

	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(url string) (*postgresRepository, error) {

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

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	// row := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)
	// if row == nil {
	// 	return nil, errors.New("account does not exist")
	// }

	// a := &Account{}
	// if err := row.Scan(&a.ID, &a.Name); err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	// rows, err := r.db.QueryContext(
	// 	ctx,
	// 	"SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",
	// 	skip,
	// 	take,
	// )

	// if err != nil {
	// 	return nil, err
	// }

	// defer rows.Close()
	// accounts := []Account{}
	// for rows.Next() {
	// 	a := &Account{}
	// 	if err := rows.Scan(&a.ID, &a.Name); err != nil {
	// 		accounts = append(accounts, *a)
	// 	}
	// }

	// if err = rows.Err(); err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
