package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/config"
)

type Router struct {
	r *gin.Engine
}

func InitRouter(
	userHandler *UserHandler,
) (*Router) {
	r := gin.New()

	pb := r.Group("/api/v1")

	// auth: public routes
	pb.POST("/register", userHandler.RegisterNewUser)

	return &Router{
		r: r,
	}
}

func (r *Router) Serve(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
