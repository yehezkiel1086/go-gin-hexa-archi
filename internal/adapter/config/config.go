package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Container struct {
	HTTP *HTTP
	DB *DB
	Token *Token
}

type HTTP struct {
	Host           string
	Port           string
	AllowedOrigins string
	AppEnv         string
}

type DB struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

type Token struct {
	Secrets  string
	Duration string
}

func Init() (*Container, error) {
	err := godotenv.Load()
	if err != nil {
		return &Container{}, err
	}

	HTTP := &HTTP{
		Host: os.Getenv("HTTP_HOST"),
		Port: os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
		AppEnv: os.Getenv("APP_ENV"),
	}

	DB := &DB{
		Name: os.Getenv("DB_NAME"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
	}
	
	Token := &Token{
		Secrets: os.Getenv("JWT_SECRET"),
		Duration: os.Getenv("TOKEN_DURATION"),
	}

	return &Container{
		HTTP: HTTP,
		DB: DB,
		Token: Token,
	}, nil
}
