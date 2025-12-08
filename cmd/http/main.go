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
		log.Error().Msg("failed to load .env configs")
		os.Exit(1)
	}
	fmt.Println(".env configs loaded successfully")

	// init context
	ctx := context.Background()

	// init postgres db
	db, err := postgres.New(ctx, conf.DB)
	if err != nil {
		log.Error().Msg("failed to initialize postgres db")
		os.Exit(1)		
	}
	fmt.Println("postgres db initialized successfully")

	// migrate dbs
	if err := db.Migrate(&domain.User{}, &domain.Product{}); err != nil {
		log.Error().Msg("failed to migrate databases")
		os.Exit(1)
	}
	fmt.Println("databases migrated successfully")

	// dependency injection
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	authSvc := service.NewAuthService(conf.JWT, userRepo)
	authHandler := handler.NewAuthHandler(conf.JWT, authSvc)

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)

	// routing
	r := handler.NewRouter(
		userHandler,
		authHandler,
		productHandler,
	)

	// serve api
	if err := r.Run(conf.HTTP); err != nil {
		log.Error().Msg("failed to serve api")
		os.Exit(1)
	}
}
