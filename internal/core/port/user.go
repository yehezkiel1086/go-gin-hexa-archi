package port

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserService interface {
	RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error)
}