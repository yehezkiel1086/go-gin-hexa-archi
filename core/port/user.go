package port

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
}

type UserService interface {
	RegisterNewUser(ctx context.Context, user *domain.User) (*domain.User, error)
	LoginUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
}
