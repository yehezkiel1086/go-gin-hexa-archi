package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
)

type UserHandler struct {
	svc port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

func (uh *UserHandler) RegisterUser(c *gin.Context) {
	var req domain.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": domain.ErrBadRequest.Error(),
		})
		return
	}

	_, err := uh.svc.RegisterUser(c, &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": domain.ErrInternal.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
	})
}

func (uh *UserHandler) GetUsers(c *gin.Context) {
	// get queries
	startStr := c.Query("start")
	endStr := c.Query("end")
	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("start and end queries are required").Error(),
		})
		return
	}

	// parse queries
	start, err := strconv.Atoi(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("invalid start query"),
	})
		return
	}
	
	end, err := strconv.Atoi(endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("invalid end query"))
		return
	}

	// get users
	users, err := uh.svc.GetUsers(c.Request.Context(), uint64(start), uint64(end))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
