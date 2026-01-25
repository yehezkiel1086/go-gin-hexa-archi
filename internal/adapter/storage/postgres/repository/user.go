package repository

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	db := ur.db.GetDB()

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
	}, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	db := ur.db.GetDB()

	var user *domain.User
	if err := db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	db := ur.db.GetDB()

	var user *domain.User
	if err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUsers(ctx context.Context, start, stop uint64) ([]domain.UserResponse, error) {
	db := ur.db.GetDB()

	var users []domain.UserResponse
	if err := db.Model(&domain.User{}).Offset(int(start)).Limit(int(stop - start + 1)).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	db := ur.db.GetDB()

	if err := db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id uint) (*domain.User, error) {
	db := ur.db.GetDB()

	var user *domain.User
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
