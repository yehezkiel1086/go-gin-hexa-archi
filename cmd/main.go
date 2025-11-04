package main

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/handler"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/storage/redis"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-employees/core/service"
)

func main() {
	// init .env configs
	conf, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ .env configs loaded successfully")

	// init context
	ctx := context.Background()

	// init postgres db
	db, err := postgres.InitDB(ctx, conf.DB)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ postgres db connected successfully")

	// init redis
	redis, err := redis.InitRedis(ctx, conf.Redis)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ redis connected successfully")
	defer redis.Close()

	// migrate db models
	err = db.Migrate(&domain.User{})
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ postgres db models migrated successfully")

	// dependency injections
	userRepo := repository.InitUserRepository(db)
	userSvc := service.InitUserService(redis, userRepo)
	userHandler := handler.InitUserHandler(userSvc)

	authSvc := service.InitAuthService(userRepo)
	authHandler := handler.InitAuthHandler(conf.JWT, authSvc)

	// init router
	r := handler.InitRouter(
		userHandler,
		authHandler,
	)

	// serve api
	if err := r.Serve(conf.HTTP); err != nil {
		panic(err)
	}
}