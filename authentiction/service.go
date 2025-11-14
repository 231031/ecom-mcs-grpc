package authentiction

import (
	"context"

	"github.com/231031/ecom-mcs-grpc/authentiction/utils"
	"github.com/segmentio/ksuid"
)

type Service interface{}
type authService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &authService{repository: r}
}

func (s *authService) PostAccountBuyer(ctx context.Context, u User) (*User, error) {
	u.ID = ksuid.New().String()

	hashed, err := utils.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashed
	// if err := s.repository.CreateAccount(ctx, a); err != nil {
	// 	return nil, err
	// }
	return &u, nil
}
