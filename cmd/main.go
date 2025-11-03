package main

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/handler"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/storage/postgres/repository"
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

	// migrate db models
	err = db.Migrate(&domain.User{})
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ postgres db models migrated successfully")

	// dependency injections
	userRepo := repository.InitUserRepository(db)
	userSvc := service.InitUserService(userRepo)
	userHandler := handler.InitUserHandler(userSvc)

	// init router
	r := handler.InitRouter(userHandler)

	// serve api
	if err := r.Serve(conf.HTTP); err != nil {
		panic(err)
	}
}