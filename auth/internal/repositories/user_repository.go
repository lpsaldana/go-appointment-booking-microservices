package repositories

import (
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByUsername(username string) (*models.User, error)
}

type userRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (u *userRepositoryImpl) CreateUser(user *models.User) error {
	return u.DB.Create(user).Error
}

func (u *userRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := u.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
