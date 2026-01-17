package port

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	GetPosts(ctx context.Context, start, end uint64) ([]domain.Post, error)
	GetPostByID(ctx context.Context, id uint) (*domain.Post, error)
	UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	DeletePost(ctx context.Context, id uint) (*domain.Post, error)
}

type PostService interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	GetPosts(ctx context.Context, start, end uint64) ([]domain.Post, error)
	GetPostByID(ctx context.Context, id uint) (*domain.Post, error)
	UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	DeletePost(ctx context.Context, id uint) (*domain.Post, error)
}
