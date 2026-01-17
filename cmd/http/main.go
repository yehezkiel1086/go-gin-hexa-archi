package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/handler"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/redis"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/service"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err.Error())
		os.Exit(1)
	}
}

func main() {
	// load .env configs
	conf, err := config.New()
	handleError(err, "unable to load .env configs")
	fmt.Println(".env configs loaded successfully")

	// init context
	ctx := context.Background()

	// init db connection
	db, err := postgres.New(ctx, conf.DB)
	handleError(err, "unable to connect with postgres db")
	fmt.Println("DB connection established successfully")

	// migrate dbs
	err = db.Migrate(&domain.User{}, &domain.Category{}, &domain.Post{})
	handleError(err, "migration failed")
	fmt.Println("DB migrated successfully")

	// init redis connection
	cache, err := redis.New(ctx, conf.Redis)
	handleError(err, "unable to connect with redis")
	defer cache.Close()

	fmt.Println("Redis connection established successfully")

	// dependency injections
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, cache)
	userHandler := handler.NewUserHandler(userSvc)

	authSvc := service.NewAuthService(conf.JWT, userRepo)
	authHandler := handler.NewAuthHandler(conf.JWT, authSvc)

	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categorySvc)

	postRepo := repository.NewPostRepository(db)
	postSvc := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postSvc)

	// init router
	r := handler.NewRouter(
		conf.HTTP,
		userHandler,
		authHandler,
		categoryHandler,
		postHandler,
	)

	// start server
	if err := r.Serve(conf.HTTP); err != nil {
		log.Fatal(err)
	}
}
