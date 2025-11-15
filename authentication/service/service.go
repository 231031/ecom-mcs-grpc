package service

import (
	"context"
	"errors"

	"github.com/231031/ecom-mcs-grpc/authentication/model"
	"github.com/231031/ecom-mcs-grpc/authentication/repository"
	"github.com/segmentio/ksuid"
)

var (
	ErrInvalidCredentials = errors.New("failed to login, invalid email or password")
)

type Service interface {
	CreateUser(ctx context.Context, u *model.User) (*model.User, error)
	LoginUser(ctx context.Context, email string, password string) (*model.TokenResponse, error)
	RefreshTokenUser(ctx context.Context, refreshToken string) (*model.TokenResponse, error)
}

type authService struct {
	repository   repository.Repository
	tokenService TokenService
}

func NewService(r repository.Repository, s TokenService) Service {
	return &authService{repository: r, tokenService: s}
}

func (s *authService) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	u.ID = ksuid.New().String()

	hashed, err := s.tokenService.hashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashed
	if err := s.repository.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *authService) LoginUser(ctx context.Context, email string, password string) (*model.TokenResponse, error) {
	u, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}

	isValid, err := s.tokenService.verifyPasswordSecure(u.Password, password)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, ErrInvalidCredentials
	}

	userAuth := &model.UserAuth{ID: u.ID, Email: u.Email, Role: u.Role}
	tokenPair, err := s.tokenService.GenerateNewPairToken(ctx, userAuth, "")
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *authService) RefreshTokenUser(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	userAuth := &model.UserAuth{}
	tokenPair, err := s.tokenService.GenerateNewPairToken(ctx, userAuth, refreshToken)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}
