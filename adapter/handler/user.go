package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/port"
)

type UserHandler struct {
	svc port.UserService
}

func InitUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

type UserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type UserResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  domain.Role `json:"role"`
}

func (h *UserHandler) RegisterNewUser(c *gin.Context) {
	// bind input
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create new user
	_, err := h.svc.RegisterNewUser(c.Request.Context(), &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     domain.EmployeeRole,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user registered successfully",
	})
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "no user in context",
		})
		return
	}

	claims, ok := user.(*domain.JWTClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "invalid user claims",
		})
		return
	}

	// get email param
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	// full information is only accessible by logged in user (other users can't access - unless admin)
	if email != claims.Email {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	// get user by email
	user, err := h.svc.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return the user
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var usersResponse []*UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, &UserResponse{
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		})
	}

	c.JSON(http.StatusOK, usersResponse)
}

// func (h *UserHandler) UpdateUser() {

// }

// func (h *UserHandler) DeleteUser() {

// }

// func (h *UserHandler) LogoutUser() {

// }