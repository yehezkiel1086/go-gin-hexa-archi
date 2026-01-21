package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
)

func AuthMiddleware(conf *config.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			// fallback: check Authorization header
			authHeader := c.GetHeader("Authorization")
			tokenString, _ = strings.CutPrefix(authHeader, "Bearer ")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": domain.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}

		// parse token
		claims, err := util.ParseToken(tokenString, conf, "access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": domain.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}

		// store user claims in context for downstream handlers
		c.Set("user", claims)
		c.Next()
	}
}

// ensures that the user has at least one of the required roles
func RoleMiddleware(requiredRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no user in context",
			})
			c.Abort()
			return
		}

		claims, ok := user.(*domain.JWTClaims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid user claims",
			})
			c.Abort()
			return
		}

		authorized := false
		for _, role := range requiredRoles {
			if claims.Role == role {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}