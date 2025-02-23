package services

import (
	"errors"

	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/repositories"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/api"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CreateUser(req *api.CreateUserRequest) (*api.CreateUserResponse, error)
	Login(req *api.LoginRequest) (*api.LoginResponse, error)
}

type authServiceImpl struct {
	Repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authServiceImpl{Repo: repo}
}

func (s *authServiceImpl) CreateUser(req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	_, err := s.Repo.FindByUsername(req.Username)
	if err == nil {
		return &api.CreateUserResponse{
			Message: "Username is not available",
			Success: false,
		}, nil
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := s.Repo.CreateUser(user); err != nil {
		return nil, err
	}

	return &api.CreateUserResponse{
		Message: "User created",
		Success: true,
	}, nil
}

func (s *authServiceImpl) Login(req *api.LoginRequest) (*api.LoginResponse, error) {
	user, err := s.Repo.FindByUsername(req.Username)

	if err != nil {
		return &api.LoginResponse{
			Token:   "",
			Success: false,
		}, errors.New("user_not_found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &api.LoginResponse{
			Token:   "",
			Success: false,
		}, errors.New("incorrect_password")
	}

	return &api.LoginResponse{
		Token:   "yet to implement",
		Success: true,
	}, nil
}
