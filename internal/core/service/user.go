package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// hash password
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPwd

	return s.repo.CreateUser(ctx, user)
}
