package port

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)	
}