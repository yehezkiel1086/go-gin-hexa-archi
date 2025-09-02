package postgres

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func Init(ctx context.Context, config *config.DB) (*DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai", config.Host, config.User, config.Password, config.Name, config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &DB{}, err
	}

	return &DB{db: db}, nil
}

func (db *DB) Migrate() error {
	if err := db.db.AutoMigrate(&domain.User{}); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetDB() *gorm.DB {
	return db.db
}
