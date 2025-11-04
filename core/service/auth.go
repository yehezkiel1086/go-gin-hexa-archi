package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/util"
)

type AuthService struct {
	repo port.UserRepository
}

func InitAuthService(repo port.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (as *AuthService) Login(ctx context.Context, user *domain.User) (*domain.User, error) {
	// check email
	userDB, err := as.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	// check password
	if err := util.ComparePassword(userDB.Password, user.Password); err != nil {
		return nil, err
	}

	return userDB, nil
}