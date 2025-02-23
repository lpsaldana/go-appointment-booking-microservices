package config

import (
	"log"

	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	DSN string
}

func NewDBConfig(dsn string) *DBConfig {
	return &DBConfig{DSN: dsn}
}

func (c *DBConfig) ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{})
	if err != nil {
		log.Printf("Error creating db connection %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Printf("Error migrating models to db %v", err)
		return nil, err
	}
	log.Println("DB connection success")
	return db, nil
}
