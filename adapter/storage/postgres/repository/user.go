package repository

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
)

type UserRepository struct {
	db *postgres.DB
}

func InitUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateNewUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	db := r.db.GetDB()

	err := db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	db := r.db.GetDB()

	var user domain.User
	err := db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	db := r.db.GetDB()

	var users []*domain.User
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
