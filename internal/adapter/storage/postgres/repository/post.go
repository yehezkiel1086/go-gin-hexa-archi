package repository

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type PostRepository struct {
	db *postgres.DB
}

func NewPostRepository(db *postgres.DB) *PostRepository {
	return &PostRepository{
		db,
	}
}

func (pr *PostRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	db := pr.db.GetDB()
	if err := db.WithContext(ctx).Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (pr *PostRepository) GetPosts(ctx context.Context) ([]domain.Post, error) {
	db := pr.db.GetDB()

	var posts []domain.Post
	if err := db.WithContext(ctx).Preload("Category").Preload("User").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (pr *PostRepository) GetPostByID(ctx context.Context, id uint) (*domain.Post, error) {
	db := pr.db.GetDB()

	var post *domain.Post
	if err := db.WithContext(ctx).Preload("Category").Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (pr *PostRepository) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	db := pr.db.GetDB()
	if err := db.WithContext(ctx).Save(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (pr *PostRepository) DeletePost(ctx context.Context, id uint) (*domain.Post, error) {
	db := pr.db.GetDB()

	var post *domain.Post
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}
