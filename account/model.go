package account

type BaseInfo struct {
	Email     string `gorm:"type:varchar(64);not null;unqiue;" json:"email"`
	Password  string `gorm:"type:varchar(255);not null;" json:"password"`
	FirstName string `gorm:"type:varchar(64);not null;" json:"first_name"`
	LastName  string `gorm:"type:varchar(64);not null;" json:"last_name"`
	Phone     string `gorm:"type:varchar(64);not null;" json:"phone"`
	Address   string `gorm:"not null; type:text;" json:"address"`
}

type Buyer struct {
	ID string `gorm:"primaryKey;type:varchar(27);" json:"id"`
	BaseInfo
}

type Seller struct {
	ID        string `gorm:"primaryKey;type:varchar(27);" json:"id"`
	StoreName string `gorm:"type:varchar(64);not null;" json:"store_name"`
	BaseInfo
}
