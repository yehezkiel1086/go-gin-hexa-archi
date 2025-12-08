package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/port"
)

type AuthHandler struct {
	conf *config.JWT
	svc port.AuthService
}

func NewAuthHandler(conf *config.JWT, svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		conf: conf,
		svc: svc,
	}
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ah.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	duration, err := strconv.Atoi(ah.conf.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("jwt_token", token, duration * 60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"jwt_token": token,
		"message": "user logged in successfully",
	})
}
