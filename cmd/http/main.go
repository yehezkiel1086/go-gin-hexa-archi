package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/yehezkiel1086/go-gin-hexa-archi/docs"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/handler"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/logger"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/storage/redis"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/service"
)

func handleError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}

//	@title			Go Gin Hexa Archi
//	@version		1.0
//	@description	Hexagonal architecture built with Go.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Yehezkiel Wiradhika
//	@contact.url	https://yehezkiel-wiradhika.vercel.app
//	@contact.email	yehezkiel1086@gmail.com

//	@host		127.0.0.1:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	// load .env configs
	conf, err := config.New()
	handleError(err, "unable to load .env configs")
	slog.Info(".env configs loaded successfully", "app", conf.App.Name, "env", conf.App.Env)

	// init logger
	logger.Set(conf.App)

	// init context
	ctx := context.Background()

	// init db connection
	db, err := postgres.New(ctx, conf.DB)
	handleError(err, "unable to connect with postgres db")
	slog.Info("postgres db connected successfully", "db", conf.DB.Host+":"+conf.DB.Port)

	// migrate dbs
	err = db.Migrate(&domain.User{}, &domain.Category{}, &domain.Post{})
	handleError(err, "migration failed")
	slog.Info("dbs migrated successfully")

	// init redis connection
	cache, err := redis.New(ctx, conf.Redis)
	handleError(err, "unable to connect with redis")
	defer cache.Close()

	slog.Info("redis connected successfully")

	// dependency injections
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, cache)
	userHandler := handler.NewUserHandler(userSvc)

	authSvc := service.NewAuthService(conf.JWT, userRepo)
	authHandler := handler.NewAuthHandler(conf.JWT, authSvc)

	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo, cache)
	categoryHandler := handler.NewCategoryHandler(categorySvc)

	postRepo := repository.NewPostRepository(db)
	postSvc := service.NewPostService(postRepo, cache)
	postHandler := handler.NewPostHandler(postSvc)

	// init router
	r := handler.NewRouter(
		conf.HTTP,
		conf.JWT,
		userHandler,
		authHandler,
		categoryHandler,
		postHandler,
	)

	// start server
	err = r.Serve()
	handleError(err, "failed to run backend server")
}
