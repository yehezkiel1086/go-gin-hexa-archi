package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App *App
		HTTP *HTTP
		DB *DB
		JWT *JWT
	}

	App struct {
		Name string
		Env  string
	}

	HTTP struct {
		Port string
		Host string
		AllowedOrigins string
	}

	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name   string
	}

	JWT struct {
		Secret string
		Duration string
	}
)

func InitConfig() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err.Error())
		}
	}

	App := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	HTTP := &HTTP{
		Port: os.Getenv("HTTP_PORT"),
		Host: os.Getenv("HTTP_HOST"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	DB := &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:   os.Getenv("DB_NAME"),
	}

	JWT := &JWT{
		Secret: os.Getenv("JWT_SECRET"),
		Duration: os.Getenv("TOKEN_DURATION"),
	}

	return &Container{
		App: App,
		HTTP: HTTP,
		DB: DB,
		JWT: JWT,
	}, nil
}