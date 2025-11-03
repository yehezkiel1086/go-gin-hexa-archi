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

func (h *UserHandler) LoginUser() {

}

func (h *UserHandler) GetUserByEmail() {

}

func (h *UserHandler) GetAllUsers() {

}

// func (h *UserHandler) UpdateUser() {

// }

// func (h *UserHandler) DeleteUser() {

// }

// func (h *UserHandler) LogoutUser() {

// }