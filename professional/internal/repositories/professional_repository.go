package repositories

import (
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/models"
	"gorm.io/gorm"
)

type ProfessionalRepository interface {
	CreateProfessional(professional *models.Professional) error
	ListProfessionals() ([]models.Professional, error)
	GetProfessionalByID(id uint) (*models.Professional, error)
}

type professionalRepositoryImpl struct {
	DB *gorm.DB
}

func NewProfessionalRepository(db *gorm.DB) ProfessionalRepository {
	return &professionalRepositoryImpl{DB: db}
}

func (r *professionalRepositoryImpl) CreateProfessional(professional *models.Professional) error {
	return r.DB.Create(professional).Error
}

func (r *professionalRepositoryImpl) ListProfessionals() ([]models.Professional, error) {
	var professionals []models.Professional
	err := r.DB.Find(&professionals).Error
	return professionals, err
}

func (r *professionalRepositoryImpl) GetProfessionalByID(id uint) (*models.Professional, error) {
	var professional models.Professional
	err := r.DB.First(&professional, id).Error
	if err != nil {
		return nil, err
	}
	return &professional, nil
}
