package main

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/handler"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/service"
)

func main() {
	// load .env configs
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ .env configs loaded successfully")

	// create context
	ctx := context.Background()

	// init db
	db, err := postgres.NewDB(ctx, conf.DB)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ DB connected successfully")

	// migrate db
	if err = db.Migrate(&domain.User{}, &domain.Category{}, &domain.Product{}); err != nil {
		panic(err)
	}
	fmt.Println("✅ DB migrated successfully")

	// dependency injection
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userhandler := handler.NewUserHandler(userSvc)

	authSvc := service.NewAuthService(conf.JWT, userRepo)
	authHandler := handler.NewAuthHandler(conf.JWT, authSvc)

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)

	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categorySvc)

	// init router
	r := handler.NewRouter(
		userhandler,
		authHandler,
		productHandler,
		categoryHandler,
	)

	// run server
	if err := r.Run(conf.HTTP); err != nil {
		panic(err)
	}
}
