package repositories

import (
	"time"

	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/models"
	"gorm.io/gorm"
)

type AgendaRepository interface {
	CreateSlot(slot *models.Slot) error
	ListAvailableSlots(professionalID uint, date time.Time) ([]models.Slot, error)
	CreateAppointment(appointment *models.Appointment) error
	UpdateSlotAvailability(slotID uint, available bool) error
	ListAppointments(clientID, professionalID uint) ([]models.Appointment, error)
	GetSlotByID(slotID uint) (*models.Slot, error) // Método añadido para obtener un slot por ID
}

type AgendaRepositoryImpl struct {
	DB *gorm.DB
}

func NewAgendaRepository(db *gorm.DB) AgendaRepository {
	return &AgendaRepositoryImpl{DB: db}
}

func (r *AgendaRepositoryImpl) CreateSlot(slot *models.Slot) error {
	return r.DB.Create(slot).Error
}

func (r *AgendaRepositoryImpl) ListAvailableSlots(professionalID uint, date time.Time) ([]models.Slot, error) {
	var slots []models.Slot
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.DB.Where("professional_id = ? AND start_time >= ? AND start_time < ? AND available = ?",
		professionalID, startOfDay, endOfDay, true).Find(&slots).Error
	return slots, err
}

func (r *AgendaRepositoryImpl) CreateAppointment(appointment *models.Appointment) error {
	return r.DB.Create(appointment).Error
}

func (r *AgendaRepositoryImpl) UpdateSlotAvailability(slotID uint, available bool) error {
	return r.DB.Model(&models.Slot{}).Where("id = ?", slotID).Update("available", available).Error
}

func (r *AgendaRepositoryImpl) ListAppointments(clientID, professionalID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	query := r.DB.Model(&models.Appointment{})
	if clientID != 0 {
		query = query.Where("client_id = ?", clientID)
	}
	if professionalID != 0 {
		query = query.Where("professional_id = ?", professionalID)
	}
	err := query.Find(&appointments).Error
	return appointments, err
}

func (r *AgendaRepositoryImpl) GetSlotByID(slotID uint) (*models.Slot, error) {
	var slot models.Slot
	err := r.DB.First(&slot, slotID).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}
