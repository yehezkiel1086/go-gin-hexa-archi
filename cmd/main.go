package main

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/handler"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/service"
)

func main() {
	// get .env configs
	config, err := config.Init()
	if err != nil {
		panic(err)
	}

	// init db
	ctx := context.Background()
	db, err := postgres.Init(ctx, config.DB)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Postgres!")
	
	// db migrations
	if err := db.Migrate(); err != nil {
		panic(err)
	}
	fmt.Println("Migration successful.")

	// dependency injections
	userRepo := repository.InitUserRepository(db)
	userService := service.InitUserService(userRepo)
	userHandler := handler.InitUserHandler(userService)

	// routing
	r, err := handler.InitRouter(
		config.HTTP,
		*userHandler,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Routes created successfully.")

	// run the server
	uri := fmt.Sprintf("%v:%v", config.HTTP.Host, config.HTTP.Port)

	fmt.Println("Server is running on", uri)
	if err := r.Serve(uri); err != nil {
		panic(err)
	}
}
