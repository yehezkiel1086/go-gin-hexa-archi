package postgres

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB(ctx context.Context, conf *config.DB) (*DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", conf.Host, conf.User, conf.Password, conf.Name, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Migrate(dbs ...any) error {
	return db.db.AutoMigrate(dbs...)
}

func (db *DB) GetDB() *gorm.DB {
	return db.db
}
