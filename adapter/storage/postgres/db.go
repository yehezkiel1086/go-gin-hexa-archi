package postgres

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-gin-hexa-employees/adapter/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func InitDB(ctx context.Context, conf *config.DB) (*DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", conf.Host, conf.User, conf.Password, conf.Name, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db = db.WithContext(ctx)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		db: db,
	}, nil
}

func (db *DB) Migrate(dbs ...any) error {
	return db.db.AutoMigrate(dbs...)
}

func (db *DB) GetDB() *gorm.DB {
	return db.db
}
