package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
)

type Router struct {
	*gin.Engine
}

func InitRouter(
	config *config.HTTP,
	userHandler UserHandler,
) (*Router, error) {
	r := gin.New()

	v1 := r.Group("/api/v1") // public routes
	v1.POST("/register", userHandler.Register)

	return &Router{r}, nil
}

func (r *Router) Serve(uri string) error {
	return r.Run(uri)
}
