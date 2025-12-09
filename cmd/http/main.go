package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/handler"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/service"
)

func main() {
	// load .env configs
	conf, err := config.New()
	if err != nil {
		errMsg := fmt.Errorf("failed to load .env configs")
		panic(errMsg)
		// log.Error().Msg(errMsg.Error())
		// os.Exit(1)
	}
	fmt.Println(".env configs loaded successfully")

	// init context
	ctx := context.Background()

	// init postgres db
	db, err := postgres.New(ctx, conf.DB)
	if err != nil {
		errMsg := fmt.Errorf("failed to initialize postgres db")
		panic(errMsg)
		// log.Error().Msg(errMsg.Error())
		// os.Exit(1)
	}
	fmt.Println("postgres db initialized successfully")

	// migrate dbs
	if err := db.Migrate(&domain.User{}, &domain.Product{}, &domain.Category{}); err != nil {
		errMsg := fmt.Errorf("failed to migrate databases")
		panic(errMsg)
		// log.Error().Msg(errMsg.Error())
		// os.Exit(1)
	}
	fmt.Println("databases migrated successfully")

	// dependency injection
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	authSvc := service.NewAuthService(conf.JWT, userRepo)
	authHandler := handler.NewAuthHandler(conf.JWT, authSvc)

	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categorySvc)

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)

	// routing
	r := handler.NewRouter(
		userHandler,
		authHandler,
		categoryHandler,
		productHandler,
	)

	// serve api
	if err := r.Run(conf.HTTP); err != nil {
		log.Error().Msg("failed to serve api")
		os.Exit(1)
	}
}
