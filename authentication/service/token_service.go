package service

import (
	"context"
	"crypto/rsa"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/231031/ecom-mcs-grpc/authentication/model"
	"github.com/231031/ecom-mcs-grpc/authentication/repository"
	"github.com/231031/ecom-mcs-grpc/authentication/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

var (
	ErrUnauth  = errors.New("the token is invalid")
	ErrExpired = errors.New("the refresh token is expired")
)

type TokenService interface {
	GenerateNewPairToken(ctx context.Context, user *model.UserAuth, prevToken string) (*model.TokenResponse, error)
	ValidateRefreshToken(refreshTokenStr string) (*model.RefreshTokenClaims, error)
	generateIDToken(user *model.UserAuth, key *rsa.PrivateKey, exp int64) (string, error)
	generateRefreshToken(id string, key string, exp int64) (*model.RefreshTokenData, error)
	HandleRefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error)
	hashPassword(password string) (string, error)
	verifyPasswordSecure(storedHash, providedPassword string) (bool, error)
}

type tokenService struct {
	AuthRepository        repository.Repository
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	TokenIDExpirationSecs int64
	RefreshExpirationSecs int64
}

func NewTokenService(repo repository.Repository, cfg *model.TokenConfig) TokenService {
	return &tokenService{
		AuthRepository:        repo,
		PrivateKey:            cfg.PrivateKey,
		PublicKey:             cfg.PublicKey,
		RefreshSecret:         cfg.RefreshSecret,
		TokenIDExpirationSecs: cfg.TokenIDExpirationSecs,
		RefreshExpirationSecs: cfg.RefreshExpirationSecs,
	}
}

func (s *tokenService) GenerateNewPairToken(ctx context.Context, userAuth *model.UserAuth, prevToken string) (*model.TokenResponse, error) {
	if prevToken != "" {
		key := fmt.Sprintf("refresh_token:%s", prevToken)
		value, err := s.AuthRepository.GetAndDelRefreshToken(ctx, key)
		if err != nil {
			return nil, ErrExpired
		}

		identical := strings.Split(value, ":")
		if len(identical) < 3 {
			return nil, ErrExpired
		}
		userAuth.ID = identical[0]
		userAuth.Email = identical[1]

		roleType, err := strconv.ParseInt(identical[2], 10, 32)
		if err != nil {
			fmt.Println("error converting string to int32:", err)
			return nil, ErrUnauth
		}
		userAuth.Role = int32(roleType)

	}

	// generate new token - login, refresh token
	newToken, err := s.generateIDToken(userAuth, s.PrivateKey, s.TokenIDExpirationSecs)
	if err != nil {
		return nil, err
	}

	newRefresh, err := s.generateRefreshToken(userAuth.ID, s.RefreshSecret, s.RefreshExpirationSecs)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("refresh_token:%s", newRefresh.SS)
	val := fmt.Sprintf("%s:%s:%s", userAuth.ID, userAuth.Email, strconv.Itoa(int(userAuth.Role)))
	err = s.AuthRepository.StoreRefreshToken(ctx, key, val, newRefresh.ExpiresIn)
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  newToken,
		RefreshToken: newRefresh.SS,
	}, nil
}

func (s *tokenService) ValidateRefreshToken(refreshTokenStr string) (*model.RefreshTokenClaims, error) {
	claims := &model.RefreshTokenClaims{}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.RefreshSecret), nil
	})
	if err != nil {
		fmt.Sprintln("failed to parse with claims refresh token", err)
		return nil, err
	}

	if !refreshToken.Valid {
		fmt.Sprintln("the refresh token is invalid with token : ", refreshToken)
		return nil, ErrUnauth
	}

	claims, ok := refreshToken.Claims.(*model.RefreshTokenClaims)
	if !ok {
		fmt.Sprintln("the token's valid but failed to parse claims")
		return nil, ErrUnauth
	}

	return claims, nil
}

func (s *tokenService) generateIDToken(user *model.UserAuth, key *rsa.PrivateKey, exp int64) (string, error) {
	curTime := time.Now()
	tokenExp := curTime.Unix() + exp

	claims := model.TokenClaims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  curTime.Unix(),
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		fmt.Sprintln("failed to signed string of token", err)
		return "", err
	}

	return ss, nil
}

func (s *tokenService) generateRefreshToken(id string, key string, exp int64) (*model.RefreshTokenData, error) {
	curTime := time.Now()
	tokenExp := curTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom()
	if err != nil {
		fmt.Sprintln("falied to generate uuid", err)
		return nil, err
	}

	claims := model.RefreshTokenClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  curTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Sprintln("failed to signed string of refresh token", err)
		return nil, err
	}

	return &model.RefreshTokenData{
		SS:        ss,
		ID:        tokenID.String(),
		ExpiresIn: tokenExp.Sub(curTime),
	}, nil
}

func (s *tokenService) HandleRefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	claims, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	_, err = uuid.Parse(claims.Id)
	if err != nil {
		fmt.Sprintln("failed to parse claims token uuid", err)
		return nil, err
	}

	userAuth := &model.UserAuth{}
	tokenPair, err := s.GenerateNewPairToken(ctx, userAuth, refreshToken)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *tokenService) hashPassword(password string) (string, error) {
	cfg := &model.Argon2Configuration{
		TimeCost:   2,
		MemoryCost: 64 * 1024,
		Threads:    4,
		KeyLength:  32,
	}

	salt, err := utils.GenerateSalt(cfg.KeyLength)
	if err != nil {
		fmt.Sprintln("failed to generate salt:", err)
		return "", err
	}
	cfg.Salt = salt

	hash := argon2.IDKey([]byte(password), cfg.Salt, cfg.TimeCost, cfg.MemoryCost, cfg.Threads, cfg.KeyLength)
	cfg.HashRaw = hash

	// Generate standardized hash format
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		cfg.MemoryCost,
		cfg.TimeCost,
		cfg.Threads,
		base64.RawStdEncoding.EncodeToString(cfg.Salt),
		base64.RawStdEncoding.EncodeToString(cfg.HashRaw),
	)

	return encodedHash, nil
}

func (s *tokenService) verifyPasswordSecure(storedHash, providedPassword string) (bool, error) {
	// Parse stored hash parameters
	config, err := utils.ParseArgon2Hash(storedHash)
	if err != nil {
		fmt.Sprintln("error parsing stored hash:", err)
		return false, err
	}

	// Generate hash using identical parameters
	computedHash := argon2.IDKey(
		[]byte(providedPassword),
		config.Salt,
		config.TimeCost,
		config.MemoryCost,
		config.Threads,
		config.KeyLength,
	)

	// Perform constant-time comparison to prevent timing attacks
	match := subtle.ConstantTimeCompare(config.HashRaw, computedHash) == 1
	return match, nil
}
