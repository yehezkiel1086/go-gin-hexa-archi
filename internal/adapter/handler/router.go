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
	httpConf *config.HTTP
	jwtConf *config.JWT
}

func NewRouter(
	httpConf *config.HTTP,
	jwtConf *config.JWT,
	userHandler *UserHandler,
	authHandler *AuthHandler,
	categoryHandler *CategoryHandler,
	postHandler *PostHandler,
) (*Router) {
	// init router
	r := gin.New()

	// define allowed origins array of string
	allowedOrigins := strings.Split(httpConf.AllowedOrigins, ",")

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
	us := pb.Group("/", AuthMiddleware(jwtConf), RoleMiddleware(domain.UserRole, domain.AdminRole))
	ad := pb.Group("/", AuthMiddleware(jwtConf), RoleMiddleware(domain.AdminRole))

	// public user and auth routes
	pb.POST("/login", authHandler.Login)
	pb.POST("/register", userHandler.RegisterUser)
	pb.GET("/refresh", authHandler.Refresh)
	pb.GET("/logout", authHandler.Logout)

	// admin user routes
	ad.GET("/users", userHandler.GetUsers)

	// public category routes
	pb.GET("/categories", categoryHandler.GetCategories)
	pb.GET("/categories/:id", categoryHandler.GetCategoryByID)

	// admin category routes
	ad.POST("/categories", categoryHandler.CreateCategory)
	ad.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	// public post routes
	pb.GET("/posts", postHandler.GetPosts)
	pb.GET("/posts/:id", postHandler.GetPostByID)

	// user post routes
	us.POST("/posts", postHandler.CreatePost)
	us.PUT("/posts/:id", postHandler.UpdatePost)
	us.DELETE("/posts/:id", postHandler.DeletePost)

	return &Router{
		r,
		httpConf,
		jwtConf,
	}
}

func (r *Router) Serve() error {
	uri := r.httpConf.Host + ":" + r.httpConf.Port
	return r.r.Run(uri)
}
