package repository

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
)

type ProductRepository struct {
	db *postgres.DB
}

func NewProductRepository(db *postgres.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	db := r.db.GetDB()

	if err := db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	db := r.db.GetDB()

	var product domain.Product
	if err := db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context) ([]domain.Product, error) {
	db := r.db.GetDB()

	var products []domain.Product
	if err := db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}