package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/config"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(
	userHandler *UserHandler,
) *Router {
	r := gin.New()

	// route groupings
	pb := r.Group("/api/v1") // public routes

	// public user routes
	pb.POST("/register", userHandler.RegisterUser)

	return &Router{r}
}

func (r *Router) Run(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
