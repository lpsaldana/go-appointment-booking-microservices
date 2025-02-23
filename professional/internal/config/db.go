package config

import (
	"log"

	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/models"
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
		log.Printf("Error connecting to DB: %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(&models.Professional{}); err != nil {
		log.Printf("Error in DB migration: %v", err)
		return nil, err
	}

	log.Println("DB conection success")
	return db, nil
}
