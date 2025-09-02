package database

import (
	"fmt"

	"kube/internal/config"
	"kube/internal/middleware"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		middleware.LogError("Failed to connect to database", err)
		panic("Database connection failed")
	}

	DB = db
	middleware.LogSuccess("Database connected successfully")
	return db
}

func GetDB() *gorm.DB {
	return DB
}
