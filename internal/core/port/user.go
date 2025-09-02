package port

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
}

type UserHandler interface {
	Register(c *gin.Context)
}
