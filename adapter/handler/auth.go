package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/port"
)

type AuthHandler struct {
	svc port.AuthService
	conf *config.JWT
}

func InitAuthHandler(conf *config.JWT, svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
		conf: conf,
	}
}

type LoginReq struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ah *AuthHandler) Login(c *gin.Context) {
	// bind json inputs
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// login user
	user, err := ah.svc.Login(c.Request.Context(), &domain.User{
		Email: req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// generate jwt token
	mySigningKey := []byte(ah.conf.Secret)

	duration, err := strconv.Atoi(ah.conf.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create claims with multiple fields populated
	claims := &domain.JWTClaims{
		Email: req.Email,
		Name: user.Name,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(duration) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	c.SetCookie("jwt_token", ss, 1000 * duration, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "user logged in successfully",
	})
}