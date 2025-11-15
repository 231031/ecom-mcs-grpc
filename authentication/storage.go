package authentication

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	return db, nil
}

func ConnectRedis(addrRedis string, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addrRedis,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return client, client.Ping(ctxTimeout).Err()
}
