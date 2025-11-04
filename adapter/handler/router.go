package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
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
	us := r.Group("/api/v1", AuthMiddleware(), RoleMiddleware(domain.EmployeeRole))
	ad := r.Group("/api/v1/admin", AuthMiddleware(), RoleMiddleware(domain.AdminRole))

	// auth: public routes
	pb.POST("/register", userHandler.RegisterNewUser)
	pb.POST("/login", authHandler.Login)

	// employees routes
	us.GET("/employees/:email", userHandler.GetUserByEmail)
	ad.GET("/employees", userHandler.GetAllUsers)

	return &Router{
		r: r,
	}
}

func (r *Router) Serve(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
