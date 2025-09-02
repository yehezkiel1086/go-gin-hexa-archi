package repository

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type UserRepository struct {
	db *postgres.DB
}

func InitUserRepository(db *postgres.DB) (*UserRepository) {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	db := ur.db.GetDB()
	if err := db.Create(user).Error; err != nil {
		return &domain.User{}, err
	}

	return user, nil
}
