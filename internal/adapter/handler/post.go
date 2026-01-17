package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
	"gorm.io/gorm"
)

type PostHandler struct {
	svc port.PostService
}

func NewPostHandler(svc port.PostService) *PostHandler {
	return &PostHandler{
		svc,
	}
}

type CreatePostReq struct {
	CategoryID  uint   `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Published   bool   `json:"published"`
}

type UpdatePostReq struct {
	CategoryID  uint   `json:"category_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Published   bool   `json:"published"`
}

func (ph *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
		return
	}

	claims, ok := user.(*domain.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternal.Error()})
		return
	}

	slug := strings.Join(strings.Split(strings.ToLower(req.Title), " "), "-")

	post, err := ph.svc.CreatePost(c, &domain.Post{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Content:     req.Content,
		Published:   req.Published,
		UserID:      claims.ID,
		Slug:        slug,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (ph *PostHandler) GetPosts(c *gin.Context) {
	// get queries
	startStr := c.Query("start")
	endStr := c.Query("end")
	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": domain.ErrInvalidQuery.Error(),
		})
		return
	}

	start, err := strconv.Atoi(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": domain.ErrInvalidQuery.Error(),
		})
		return
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": domain.ErrInvalidQuery.Error(),
		})
		return
	}

	// get posts
	posts, err := ph.svc.GetPosts(c, uint64(start), uint64(end))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": domain.ErrInternal.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": domain.ErrBadRequest,
		})
		return
	}

	post, err := ph.svc.GetPostByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (ph *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req UpdatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post *domain.Post

	post, err = ph.svc.UpdatePost(c, &domain.Post{
		Model:       gorm.Model{ID: uint(id)},
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Content:     req.Content,
		Published:   req.Published,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (ph *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	post, err := ph.svc.DeletePost(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}
