package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/storage/redis"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/util"
)

type UserService struct {
	repo port.UserRepository
	redis *redis.Redis
}

func InitUserService(redis *redis.Redis, repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
		redis: redis,
	}
}

func (us *UserService) RegisterNewUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPwd

	dbUser, err := us.repo.CreateNewUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return dbUser, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return us.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.GetAllUsers(ctx)
}