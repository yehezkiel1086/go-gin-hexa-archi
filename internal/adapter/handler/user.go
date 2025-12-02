package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/port"
)

type UserHandler struct {
	svc port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{svc}
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name string `json:"name" binding:"required"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name: req.Name,
	}

	_, err := h.svc.RegisterUser(c, user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
	})
}
