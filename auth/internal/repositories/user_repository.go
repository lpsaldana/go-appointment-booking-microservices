package repositories

import "github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/models"

type UserRepository interface {
	ValidateCredentials(string, string) (*models.User, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (u *userRepository) ValidateCredentials(username string, password string) (*models.User, error) {
	return nil, nil
}
