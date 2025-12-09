package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/port"
)

type CategoryService struct {
	repo port.CategoryRepository
}

func NewCategoryService(repo port.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	return s.repo.CreateCategory(ctx, category)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetCategories(ctx)
}
