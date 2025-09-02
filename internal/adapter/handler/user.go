package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
)

type UserHandler struct {
	repo port.UserService
}

func InitUserHandler(repo port.UserService) (*UserHandler) {
	return &UserHandler{
		repo: repo,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func (uh *UserHandler) Register(c *gin.Context) {
	// insert user inputs
	var input RegisterRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// create new user
	user := &domain.User{
		Username: input.Username,
		Password: input.Password,
		Fullname: input.Fullname,
		Email: input.Email,
	}

	// register new user
	_, err := uh.repo.Register(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User is registered successfully.",
	})
}
