package main

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
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
}
