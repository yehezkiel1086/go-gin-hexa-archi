package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/util"
)

type UserService struct {
	repo port.UserRepository
}

func InitUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) RegisterNewUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPwd

	return us.repo.CreateNewUser(ctx, user)
}

func (us *UserService) LoginUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// compare email
	userDB, err := us.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	// compare password
	err = util.ComparePassword(user.Password, userDB.Password)
	if err != nil {
		return nil, err
	}

	return userDB, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return us.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.GetAllUsers(ctx)
}