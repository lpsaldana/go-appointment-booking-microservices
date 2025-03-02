package unit

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	secretKey := "test-secret-key"
	mockRepo := new(MockUserRepository)
	srv := services.NewAuthService(mockRepo, secretKey)

	tests := []struct {
		name         string
		req          *pb.CreateUserRequest
		mockSetup    func()
		expectedResp *pb.CreateUserResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.CreateUserRequest{Username: "testuser", Password: "testpass"},
			mockSetup: func() {
				mockRepo.On("FindByUsername", "testuser").Return((*models.User)(nil), errors.New("not found")).Once()
				mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil).Once()
			},
			expectedResp: &pb.CreateUserResponse{Message: "User created", Success: true},
			expectedErr:  nil,
		},
		{
			name: "UserAlreadyExists",
			req:  &pb.CreateUserRequest{Username: "existinguser", Password: "testpass"},
			mockSetup: func() {
				mockRepo.On("FindByUsername", "existinguser").Return(&models.User{Username: "existinguser"}, nil).Once()
			},
			expectedResp: &pb.CreateUserResponse{Message: "Username is not available", Success: false},
			expectedErr:  nil,
		},
		{
			name: "DatabaseError",
			req:  &pb.CreateUserRequest{Username: "testuser", Password: "testpass"},
			mockSetup: func() {
				mockRepo.On("FindByUsername", "testuser").Return((*models.User)(nil), errors.New("not found")).Once()
				mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(errors.New("db error")).Once()
			},
			expectedResp: nil,
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.CreateUser(tt.req)
			assert.Equal(t, tt.expectedResp, resp)
			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	secretKey := "test-secret-key"
	mockRepo := new(MockUserRepository)
	srv := services.NewAuthService(mockRepo, secretKey)

	// Mock de usuario con contraseña encriptada
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	user := &models.User{ID: 1, Username: "testuser", Password: string(hashedPass)}

	tests := []struct {
		name         string
		req          *pb.LoginRequest
		mockSetup    func()
		expectedResp *pb.LoginResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.LoginRequest{Username: "testuser", Password: "testpass"},
			mockSetup: func() {
				mockRepo.On("FindByUsername", "testuser").Return(user, nil).Once()
			},
			expectedResp: &pb.LoginResponse{Success: true},
			expectedErr:  nil,
		},
		{
			name: "UserNotFound",
			req:  &pb.LoginRequest{Username: "unknown", Password: "testpass"},
			mockSetup: func() {
				mockRepo.On("FindByUsername", "unknown").Return((*models.User)(nil), errors.New("not found")).Once()
			},
			expectedResp: &pb.LoginResponse{Token: "", Success: false},
			expectedErr:  errors.New("user_not_found"),
		},
		{
			name: "WrongPassword",
			req:  &pb.LoginRequest{Username: "testuser", Password: "wrongpass"},
			mockSetup: func() {
				mockRepo.On("FindByUsername", "testuser").Return(user, nil).Once()
			},
			expectedResp: &pb.LoginResponse{Token: "", Success: false},
			expectedErr:  errors.New("incorrect_password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.Login(tt.req)
			assert.Equal(t, tt.expectedResp.Success, resp.Success)
			if tt.expectedResp.Success {
				assert.NotEmpty(t, resp.Token)
				token, _ := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
					return []byte(secretKey), nil
				})
				assert.True(t, token.Valid)
			}
			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
