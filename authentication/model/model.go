package model

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"password"`
	Role     int32  `json:"role"`
}

type UserAuth struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Email string `json:"email" gorm:"uniqueIndex"`
	Role  int32  `json:"role"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Argon2Configuration struct {
	HashRaw    []byte
	Salt       []byte
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
	KeyLength  uint32
}

type TokenConfig struct {
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	TokenIDExpirationSecs int64
	RefreshExpirationSecs int64
}

type TokenClaims struct {
	User UserAuth `json:"user"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type RefreshTokenData struct {
	SS        string
	ID        string
	ExpiresIn time.Duration
}
