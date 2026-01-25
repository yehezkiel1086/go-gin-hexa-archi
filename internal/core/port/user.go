package port

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

//go:generate mockery --name=UserRepository --output=../../../mocks --outpkg=mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUsers(ctx context.Context, start, stop uint64) ([]domain.UserResponse, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint) (*domain.User, error)
}

//go:generate mockery --name=UserService --output=../../../mocks --outpkg=mocks
type UserService interface {
	RegisterUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error)
	GetUsers(ctx context.Context, start, stop uint64) ([]domain.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, id uint, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint) (*domain.User, error)
}
