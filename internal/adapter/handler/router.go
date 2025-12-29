package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(
	userHandler *UserHandler,
	authHandler *AuthHandler,
) *Router {
	r := gin.New()

	// route groupings
	pb := r.Group("/api/v1")

	// public user routes
	pb.POST("/register", userHandler.RegisterUser)
	pb.POST("/login", authHandler.Login)

	return &Router{r}
}

func (r *Router) Run(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
