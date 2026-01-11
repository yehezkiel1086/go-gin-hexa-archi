package repository

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type CategoryRepository struct {
	db *postgres.DB
}

func NewCategoryRepository(db *postgres.DB) *CategoryRepository {
	return &CategoryRepository{
		db,
	}
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	db := cr.db.GetDB()
	if err := db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (cr *CategoryRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	db := cr.db.GetDB()

	var categories []domain.Category
	if err := db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (cr *CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (*domain.Category, error) {
	db := cr.db.GetDB()

	var category *domain.Category
	if err := db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (cr *CategoryRepository) DeleteCategory(ctx context.Context, id uint) (*domain.Category, error) {
	db := cr.db.GetDB()

	var category *domain.Category
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&category).Error; err != nil {
		return nil, err
	}

	return category, nil
}
