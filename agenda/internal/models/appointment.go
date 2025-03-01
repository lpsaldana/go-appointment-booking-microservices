package models

type Appointment struct {
	ID             uint `gorm:"primaryKey"`
	ClientID       uint `gorm:"not null"`
	SlotID         uint `gorm:"not null;unique"`
	ProfessionalID uint `gorm:"not null"`
}
