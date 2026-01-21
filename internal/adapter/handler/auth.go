package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
)

type AuthHandler struct {
	conf *config.JWT
	svc port.AuthService
}

func NewAuthHandler(conf *config.JWT, svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		conf,
		svc,
	}
}

type LoginUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req LoginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("email and password are required")})
		return
	}

	refreshToken, accessToken, err := ah.svc.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized})
		return
	}

	// convert duration to int
	refreshTokenDuration, err := strconv.Atoi(ah.conf.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": domain.ErrInternal,
		})
		return
	}

	accessTokenDuration, err := strconv.Atoi(ah.conf.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": domain.ErrInternal,
		})
		return
	}

	// set jwt token in cookie
	c.SetCookie("refresh_token", refreshToken, refreshTokenDuration * 60 * 60 * 24, "/api/v1/refresh", "", false, true)
	c.SetCookie("access_token", accessToken, accessTokenDuration, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func (ah *AuthHandler) Refresh(c *gin.Context) {
	// get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate refresh token
	claims, err := util.ParseToken(refreshToken, ah.conf, "refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := ah.svc.GetUserByEmail(c.Request.Context(), claims.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// generate new access token
	accessToken, err := util.GenerateJWTToken(ah.conf, user, "access")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// set to cookie
	duration, err := strconv.Atoi(ah.conf.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", accessToken, duration, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"access_token": accessToken,
	})
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "user logged out",
	})
}
