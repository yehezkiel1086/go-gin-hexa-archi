package repository

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
)

type CategoryRepository struct {
	db *postgres.DB
}

func NewCategoryRepository(db *postgres.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	db := r.db.GetDB()

	if err := db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	db := r.db.GetDB()

	var category domain.Category

	if err := db.First(&category, uint(id)).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	db := r.db.GetDB()

	var categories []domain.Category

	if err := db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
