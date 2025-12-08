package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/port"
)

type ProductService struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) GetProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.GetProducts(ctx)
}