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
	authHandler *AuthHandler,
) (*Router) {
	r := gin.New()

	pb := r.Group("/api/v1")
	ad := r.Group("/api/v1/admin")

	// auth: public routes
	pb.POST("/register", userHandler.RegisterNewUser)
	pb.POST("/login", authHandler.Login)

	// employees routes
	ad.GET("/employees", userHandler.GetAllUsers)
	ad.GET("/employees/:email", userHandler.GetUserByEmail)

	return &Router{
		r: r,
	}
}

func (r *Router) Serve(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
