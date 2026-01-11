package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
)

type CategoryService struct {
	repo port.CategoryRepository
}

func NewCategoryService(repo port.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	return cs.repo.CreateCategory(ctx, category)
}

func (cs *CategoryService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return cs.repo.GetCategories(ctx)
}

func (cs *CategoryService) GetCategoryByID(ctx context.Context, id uint) (*domain.Category, error) {
	return cs.repo.GetCategoryByID(ctx, id)
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, id uint) (*domain.Category, error) {
	return cs.repo.DeleteCategory(ctx, id)
}
