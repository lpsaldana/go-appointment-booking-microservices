package services

import (
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/repositories"
)

type UserService interface {
	ValidateCredentials(string, string) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) ValidateCredentials(username string, password string) (*models.User, error) {
	return s.userRepository.ValidateCredentials(username, password)
}
