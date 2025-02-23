package services

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	Repo      repositories.UserRepository
	SecretKey []byte
}

func NewAuthService(repo repositories.UserRepository, secretKey string) AuthService {
	return &authServiceImpl{Repo: repo, SecretKey: []byte(secretKey)}
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

	log.Println(user.Password, req.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &api.LoginResponse{
			Token:   "",
			Success: false,
		}, errors.New("incorrect_password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,                               // ID del usuario en el cuerpo
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Expira en 24 horas
		"iat":     time.Now().Unix(),                     // Issued At: tiempo de emisión
	})

	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		log.Printf("Error singning token: %v", err)
		return &api.LoginResponse{
			Token:   "",
			Success: false,
		}, errors.New("error_generating_token")
	}

	return &api.LoginResponse{
		Token:   tokenString,
		Success: true,
	}, nil
}
