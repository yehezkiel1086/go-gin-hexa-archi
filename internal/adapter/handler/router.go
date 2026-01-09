package handler

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(
	conf *config.HTTP,
	userHandler *UserHandler,
	authHandler *AuthHandler,
) (*Router) {
	// init router
	r := gin.New()

	// define allowed origins array of string
	allowedOrigins := strings.Split(conf.AllowedOrigins, ",")

	// cors config
	corsConf := cors.New(cors.Config{
    AllowOrigins:   allowedOrigins,
    AllowMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
	})
	r.Use(corsConf)

	// group routes
	pb := r.Group("/api/v1")
	// us := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.UserRole, domain.AdminRole))
	ad := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.AdminRole))

	// public user and auth routes
	pb.POST("/login", authHandler.Login)
	pb.POST("/register", userHandler.RegisterUser)

	// admin user routes
	ad.GET("/users", userHandler.GetUsers)

	return &Router{r}
}

func (r *Router) Serve(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
