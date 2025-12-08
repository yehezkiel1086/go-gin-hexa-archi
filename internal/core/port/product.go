package port

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductByID(ctx context.Context, id int) (*domain.Product, error)
	GetProducts(ctx context.Context) ([]domain.Product, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductByID(ctx context.Context, id int) (*domain.Product, error)
	GetProducts(ctx context.Context) ([]domain.Product, error)	
}