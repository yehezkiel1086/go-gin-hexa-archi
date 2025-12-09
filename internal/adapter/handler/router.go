package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(
	userHandler *UserHandler,
	authHandler *AuthHandler,
	categoryHandler *CategoryHandler,
	productsHandler *ProductHandler,
) *Router {
	r := gin.New()

	// route groupings
	pb := r.Group("/api/v1") // public routes
	us := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.UserRole, domain.AdminRole))
	ad := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.AdminRole))

	// public user routes
	pb.POST("/register", userHandler.RegisterUser)
	pb.POST("/login", authHandler.Login)

	// user category routes
	us.GET("/categories", categoryHandler.GetCategories)
	us.GET("/categories/:id", categoryHandler.GetCategoryByID)

	// admin category routes
	ad.POST("/categories", categoryHandler.CreateCategory)

	// user product routes
	us.GET("/products", productsHandler.GetProducts)
	us.GET("/products/:id", productsHandler.GetProductByID)

	// admin product routes
	ad.POST("/products", productsHandler.CreateProduct)

	return &Router{r}
}

func (r *Router) Run(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
