package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
)

type CategoryService struct {
	repo  port.CategoryRepository
	cache port.CacheRepository
}

func NewCategoryService(repo port.CategoryRepository, cache port.CacheRepository) *CategoryService {
	return &CategoryService{
		repo,
		cache,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	category, err := cs.repo.CreateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	cacheKey := util.GenerateCacheKey("category", category.ID)
	categorySerialized, err := util.Serialize(category)
	if err != nil {
		return nil, err
	}

	if err := cs.cache.Set(ctx, cacheKey, categorySerialized, 0); err != nil {
		return nil, err
	}

	if err := cs.cache.Delete(ctx, "categories"); err != nil {
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	cacheKey := "categories"

	categoriesSerialized, err := cs.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(categoriesSerialized, &categories); err != nil {
			return nil, err
		}
		return categories, nil
	}

	categories, err = cs.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	categoriesSerialized, err = util.Serialize(categories)
	if err != nil {
		return nil, err
	}

	if err := cs.cache.Set(ctx, cacheKey, categoriesSerialized, 0); err != nil {
		return nil, err
	}

	return categories, nil
}

func (cs *CategoryService) GetCategoryByID(ctx context.Context, id uint) (*domain.Category, error) {
	category := &domain.Category{}
	cacheKey := util.GenerateCacheKey("category", id)

	categorySerialized, err := cs.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(categorySerialized, category); err != nil {
			return nil, err
		}
		return category, nil
	}

	category, err = cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	categorySerialized, err = util.Serialize(category)
	if err != nil {
		return nil, err
	}

	if err := cs.cache.Set(ctx, cacheKey, categorySerialized, 0); err != nil {
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, id uint) (*domain.Category, error) {
	cacheKey := util.GenerateCacheKey("category", id)

	if err := cs.cache.Delete(ctx, cacheKey); err != nil {
		return nil, err
	}

	if err := cs.cache.Delete(ctx, "categories"); err != nil {
		return nil, err
	}

	return cs.repo.DeleteCategory(ctx, id)
}
