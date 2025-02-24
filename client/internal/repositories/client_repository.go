package repositories

import (
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/models"
	"gorm.io/gorm"
)

type ClientRepository interface {
	CreateClient(client *models.Client) error
	GetClientByID(id uint) (*models.Client, error)
	ListClients() ([]models.Client, error)
}

type ClientRepositoryImpl struct {
	DB *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &ClientRepositoryImpl{DB: db}
}

func (r *ClientRepositoryImpl) CreateClient(client *models.Client) error {
	return r.DB.Create(client).Error
}

func (r *ClientRepositoryImpl) GetClientByID(id uint) (*models.Client, error) {
	var client models.Client
	err := r.DB.First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepositoryImpl) ListClients() ([]models.Client, error) {
	var clients []models.Client
	err := r.DB.Find(&clients).Error
	return clients, err
}
