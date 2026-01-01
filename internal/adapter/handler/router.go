package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(
	userHandler *UserHandler,
	authHandler *AuthHandler,
	productHandler *ProductHandler,
	categoryHandler *CategoryHandler,
) *Router {
	r := gin.New()

	// route groupings
	pb := r.Group("/api/v1")
	us := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.UserRole, domain.AdminRole))
	ad := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.AdminRole))

	// public user routes
	pb.POST("/register", userHandler.RegisterUser)
	pb.POST("/login", authHandler.Login)

	// public product routes
	pb.GET("/products", productHandler.GetProducts)
	pb.GET("/products/:id", productHandler.GetProductByID)

	// admin product routes
	ad.POST("/products", productHandler.CreateProduct)
	ad.DELETE("/products/:id", productHandler.DeleteProduct)

	// user category routes
	us.GET("/categories", categoryHandler.GetCategories)
	us.GET("/categories/:id", categoryHandler.GetCategoryByID)
	
	// admin category routes
	ad.POST("/categories", categoryHandler.CreateCategory)
	ad.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	return &Router{r}
}

func (r *Router) Run(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
