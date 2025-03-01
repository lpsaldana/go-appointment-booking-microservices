package models

import "time"

type Slot struct {
	ID             uint      `gorm:"primaryKey"`
	ProfessionalID uint      `gorm:"not null"`
	StartTime      time.Time `gorm:"not null"`
	EndTime        time.Time `gorm:"not null"`
	Available      bool      `gorm:"default:true"`
}
