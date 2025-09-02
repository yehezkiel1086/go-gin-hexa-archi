package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
}

func InitUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return &domain.User{}, err
	}

	user.Password = hashedPwd

	user, err = us.repo.CreateUser(ctx, user)
	if err != nil {
		return &domain.User{}, err
	}

	return user, nil
}
