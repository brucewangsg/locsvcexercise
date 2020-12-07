package coresvc

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type appConfig struct {
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
}

func NewAppConfig() *appConfig {
	config := new(appConfig)
	config.DatabaseHost = os.Getenv("DB_HOST")
	config.DatabasePort = os.Getenv("DB_PORT")
	config.DatabaseUser = os.Getenv("DB_USER")
	config.DatabasePass = os.Getenv("DB_PASS")
	config.DatabaseName = os.Getenv("DB_NAME")

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

	return config
}

func NewAppDBPool(config *appConfig) *gorm.DB {
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
		Logger: logger.Default.LogMode(logger.Silent),
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
