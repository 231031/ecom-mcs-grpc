package account

import "time"

type BaseInfo struct {
	FirstName string `gorm:"type:varchar(64);not null;" json:"first_name"`
	LastName  string `gorm:"type:varchar(64);not null;" json:"last_name"`
	Phone     string `gorm:"type:varchar(64);not null;" json:"phone"`
	Address   string `gorm:"not null; type:text;" json:"address"`
}

type Buyer struct {
	ID string `gorm:"primaryKey;type:varchar(27);" json:"id"`
	BaseInfo
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Seller struct {
	ID        string `gorm:"primaryKey;type:varchar(27);" json:"id"`
	StoreName string `gorm:"type:varchar(64);not null;" json:"store_name"`
	BaseInfo
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
