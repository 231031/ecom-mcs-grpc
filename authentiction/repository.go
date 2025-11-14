package authentiction

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	Close() error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(url string) (Repository, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	return &repository{db}, nil
}

func (r *repository) Close() error {
	postgres, err := r.db.DB()
	if err != nil {
		return err
	}
	return postgres.Close()
}
