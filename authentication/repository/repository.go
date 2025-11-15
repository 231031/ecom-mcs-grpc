package repository

import (
	"context"
	"time"

	"github.com/231031/ecom-mcs-grpc/authentication/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	Close() error
	CreateUser(ctx context.Context, u *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	StoreRefreshToken(ctx context.Context, key string, value string, exp time.Duration) error
	GetAndDelRefreshToken(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, key string) error
}

type repository struct {
	db          *gorm.DB
	redisClient *redis.Client
	ctx         context.Context
}

func NewRepository(db *gorm.DB, redisClient *redis.Client) Repository {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return &repository{db: db, redisClient: redisClient, ctx: ctx}
}

func (r *repository) Close() error {
	postgres, err := r.db.DB()
	if err != nil {
		return err
	}
	return postgres.Close()
}

func (r *repository) CreateUser(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// Refresh Token operations
func (r *repository) StoreRefreshToken(ctx context.Context, key string, value string, exp time.Duration) error {
	if err := r.redisClient.Set(ctx, key, value, exp).Err(); err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAndDelRefreshToken(ctx context.Context, key string) (string, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *repository) DeleteRefreshToken(ctx context.Context, key string) error {
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
