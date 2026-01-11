package service

import (
	"context"
	"strings"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
)

type PostService struct {
	repo port.PostRepository
}

func NewPostService(repo port.PostRepository) *PostService {
	return &PostService{
		repo,
	}
}

func (ps *PostService) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	return ps.repo.CreatePost(ctx, post)
}

func (ps *PostService) GetPosts(ctx context.Context) ([]domain.Post, error) {
	return ps.repo.GetPosts(ctx)
}

func (ps *PostService) GetPostByID(ctx context.Context, id uint) (*domain.Post, error) {
	return ps.repo.GetPostByID(ctx, id)
}

func (ps *PostService) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	found, err := ps.repo.GetPostByID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	if post.Title != "" {
		found.Title = post.Title
		slug := strings.Join(strings.Split(strings.ToLower(post.Title), " "), "-")
		found.Slug = slug
	}
	if post.Content != "" {
		found.Content = post.Content
	}
	if post.CategoryID != 0 {
		found.CategoryID = post.CategoryID
	}
	if post.Published != false {
		found.Published = post.Published
	}

	return ps.repo.UpdatePost(ctx, found)
}

func (ps *PostService) DeletePost(ctx context.Context, id uint) (*domain.Post, error) {
	return ps.repo.DeletePost(ctx, id)
}
