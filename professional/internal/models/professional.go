package models

type Professional struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Profession string `gorm:"not null"`
	Contact    string `gorm:"not null"`
}
