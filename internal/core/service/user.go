package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
)

type UserService struct {
	repo  port.UserRepository
	cache port.CacheRepository
}

func NewUserService(repo port.UserRepository, cache port.CacheRepository) *UserService {
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

func (us *UserService) GetUsers(ctx context.Context, start, end uint64) ([]domain.UserResponse, error) {
	users := []domain.UserResponse{}

	// get from cache
	params := util.GenerateCacheKeyParams(start, end)
	cacheKey := util.GenerateCacheKey("users", params)
	usersCache, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(usersCache, &users); err != nil {
			return nil, err
		}
		return users, nil
	}

	// get from db if doesn't exist in cache
	users, err = us.repo.GetUsers(ctx, start, end)
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

func (us *UserService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	user := &domain.User{}

	// get from cache
	cacheKey := util.GenerateCacheKey("user", id)

	serialized, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(serialized, user); err != nil {
			return nil, err
		}

		return user, nil
	}

	// get from db if not in cache
	user, err = us.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// cache user
	serialized, err = util.Serialize(user)
	if err != nil {
		return nil, err
	}

	if err := us.cache.Set(ctx, cacheKey, serialized, 0); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) UpdateUser(ctx context.Context, id uint, user *domain.User) (*domain.User, error) {
	foundUser := &domain.User{}

	// get user from cache
	cacheKey := util.GenerateCacheKey("user", id)

	serialized, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(serialized, foundUser); err != nil {
			return nil, err
		}
	}

	// get user from db if not in cache
	emptyUser := &domain.User{}
	if foundUser == emptyUser {
		foundUser, err = us.repo.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	// update user
	if user.Name != "" {
		foundUser.Name = user.Name
	}
	if user.Email != "" {
		foundUser.Email = user.Email
	}
	if user.Password != "" {
		hashed, err := util.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}

		foundUser.Password = hashed
	}

	// update user
	if _, err := us.repo.UpdateUser(ctx, foundUser); err != nil {
		return nil, err
	}

	// delete caches
	if err := us.cache.Delete(ctx, cacheKey); err != nil {
		return nil, err
	}

	if err := us.cache.DeleteByPrefix(ctx, "users:*"); err != nil {
		return nil, err
	}

	// cache updated user
	serialized, err = util.Serialize(foundUser)
	if err != nil {
		return nil, err
	}

	if err := us.cache.Set(ctx, cacheKey, serialized, 0); err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id uint) (*domain.User, error) {
	// delete user from cache
	cacheKey := util.GenerateCacheKey("user", id)
	if err := us.cache.Delete(ctx, cacheKey); err != nil {
		return nil, err
	}

	if err := us.cache.DeleteByPrefix(ctx, "users:*"); err != nil {
		return nil, err
	}

	// delete user from db
	return us.repo.DeleteUser(ctx, id)
}
