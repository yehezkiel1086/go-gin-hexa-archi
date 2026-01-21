package port

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
