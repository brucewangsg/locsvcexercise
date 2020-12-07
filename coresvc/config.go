package coresvc

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AppConfig database and app port
type AppConfig struct {
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	AppPort      string
}

// NewAppConfig get app config based on environment variable
func NewAppConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	config := new(AppConfig)
	config.DatabaseHost = os.Getenv("DB_HOST")
	config.DatabasePort = os.Getenv("DB_PORT")
	config.DatabaseUser = os.Getenv("DB_USER")
	config.DatabasePass = os.Getenv("DB_PASS")
	config.DatabaseName = os.Getenv("DB_NAME")
	config.AppPort = os.Getenv("APP_PORT")

	if config.DatabaseHost == "" {
		config.DatabaseHost = "db"
	}

	if config.DatabasePort == "" {
		config.DatabasePort = "5432"
	}

	if config.DatabaseUser == "" {
		config.DatabaseUser = "postgres"
	}

	if config.DatabaseName == "" {
		config.DatabaseName = "locexercise"
	}

	if config.AppPort == "" {
		config.AppPort = "5678"
	}

	return config
}

// NewAppDBPool get db pool instance
func NewAppDBPool(config *AppConfig) *gorm.DB {
	dbstr := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabasePass,
		config.DatabaseName,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbstr,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(4)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		panic("Something is wrong with database")
	}

	return db
}
