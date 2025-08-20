package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func GetDBEnv() string {
  // get .env values of DB
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	DB_SSL_MODE := os.Getenv("DB_SSL_MODE")

  dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=Asia/Jakarta", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSL_MODE)

  return dsn
}
