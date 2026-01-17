package service

import (
	"context"
	"strings"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/redis"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
)

type PostService struct {
	repo port.PostRepository
	cache port.CacheRepository
}

func NewPostService(repo port.PostRepository, cache *redis.Redis) *PostService {
	return &PostService{
		repo,
		cache,
	}
}

func (ps *PostService) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	// create post
	post, err := ps.repo.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	// set cache
	cacheKey := util.GenerateCacheKey("post", post.ID)
	postSerialized, err := util.Serialize(post)
	if err != nil {
		return nil, err
	}
	
	if err := ps.cache.Set(ctx, cacheKey, postSerialized, 0); err != nil {
		return nil, err
	}

	// clear posts cache (since new post created)
	if err := ps.cache.DeleteByPrefix(ctx, "posts"); err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *PostService) GetPosts(ctx context.Context, start, end uint64) ([]domain.Post, error) {
	var posts []domain.Post

	// generate cache key
	param := util.GenerateCacheKeyParams(start, end)
	cacheKey := util.GenerateCacheKey("posts", param)

	// get from cache
	postSerialized, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(postSerialized, &posts); err != nil {
			return nil, err
		}

		return posts, nil
	}

	// get from db if cache don't exist
	posts, err = ps.repo.GetPosts(ctx, start, end)
	if err != nil {
		return nil, err
	}

	// set cache
	postsSerialized, err := util.Serialize(posts)
	if err != nil {
		return nil, err
	}

	if err := ps.cache.Set(ctx, cacheKey, postsSerialized, 0); err != nil {
		return nil, err
	}

	return posts, nil
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
