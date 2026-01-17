package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
	cache port.CacheRepository
}

func NewUserService(repo port.UserRepository, cache port.CacheRepository) (*UserService) {
	return &UserService{
		repo,
		cache,
	}
}

func (us *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	// hash password
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPwd

	// set cache
	cacheKey := util.GenerateCacheKey("user", user.ID)
	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, err
	}

	if err := us.cache.Set(ctx, cacheKey, userSerialized, 0); err != nil {
		return nil, err
	}

	// clear users cache (since new user created)
	if err := us.cache.DeleteByPrefix(ctx, "users:*"); err != nil {
		return nil, err
	}

	return us.repo.CreateUser(ctx, user)
}

func (us *UserService) GetUsers(ctx context.Context, start, stop uint64) ([]domain.UserResponse, error) {
	var users []domain.UserResponse

	// get from cache
	params := util.GenerateCacheKeyParams(start, stop)
	cacheKey := util.GenerateCacheKey("users", params)
	usersCache, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(usersCache, users); err != nil {
			return nil, err
		}
		return users, nil
	}

	// get from db if doesn't exist in cache
	users, err = us.repo.GetUsers(ctx, start, stop)
	if err != nil {
		return nil, err
	}

	// cache gotten users from db
	serializedUsers, err := util.Serialize(users)
	if err != nil {
		return nil, err
	}

	if err := us.cache.Set(ctx, cacheKey, serializedUsers, 0); err != nil {
		return nil, err
	}

	return users, nil
}
