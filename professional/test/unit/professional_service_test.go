package unit

import (
	"errors"
	"testing"

	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProfessionalRepository struct {
	mock.Mock
}

func (m *MockProfessionalRepository) CreateProfessional(professional *models.Professional) error {
	args := m.Called(professional)
	return args.Error(0)
}

func (m *MockProfessionalRepository) GetProfessionalByID(id uint) (*models.Professional, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Professional), args.Error(1)
}

func (m *MockProfessionalRepository) ListProfessionals() ([]models.Professional, error) {
	args := m.Called()
	return args.Get(0).([]models.Professional), args.Error(1)
}

func TestCreateProfessionalService(t *testing.T) {
	mockRepo := new(MockProfessionalRepository)
	srv := services.NewProfessionalService(mockRepo)

	tests := []struct {
		name         string
		req          *pb.CreateProfessionalRequest
		mockSetup    func()
		expectedResp *pb.CreateProfessionalResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.CreateProfessionalRequest{Name: "Dr. Lopez", Profession: "Dentista", Contact: "lopez@email.com"},
			mockSetup: func() {
				(mockRepo).On("CreateProfessional", mock.AnythingOfType("*models.Professional")).Return(nil).Once()
			},
			expectedResp: &pb.CreateProfessionalResponse{Message: "Professional created", Success: true, ProfessionalId: 0},
			expectedErr:  nil,
		},
		{
			name: "DatabaseError",
			req:  &pb.CreateProfessionalRequest{Name: "Dr. Perez", Profession: "Medico", Contact: "perez@email.com"},
			mockSetup: func() {
				(mockRepo).On("CreateProfessional", mock.AnythingOfType("*models.Professional")).Return(errors.New("db error")).Once()
			},
			expectedResp: &pb.CreateProfessionalResponse{Message: "Error creating professional", Success: false},
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.CreateProfessional(tt.req)
			assert.Equal(t, tt.expectedResp.Message, resp.Message)
			assert.Equal(t, tt.expectedResp.Success, resp.Success)
			/*if tt.expectedResp.Success {
				assert.NotZero(t, resp.ProfessionalId, "ProfessionalId debería asignarse")
			}*/
			assert.Equal(t, tt.expectedErr, err)
			(mockRepo).AssertExpectations(t)
		})
	}
}

func TestGetProfessional(t *testing.T) {
	mockRepo := new(MockProfessionalRepository)
	srv := services.NewProfessionalService(mockRepo)

	tests := []struct {
		name         string
		req          *pb.GetProfessionalRequest
		mockSetup    func()
		expectedResp *pb.GetProfessionalResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.GetProfessionalRequest{Id: 1},
			mockSetup: func() {
				(mockRepo).On("GetProfessionalByID", uint(1)).Return(&models.Professional{ID: 1, Name: "Dr. Lopez", Profession: "Dentista", Contact: "lopez@email.com"}, nil).Once()
			},
			expectedResp: &pb.GetProfessionalResponse{
				Professional: &pb.Professional{Id: 1, Name: "Dr. Lopez", Profession: "Dentista", Contact: "lopez@email.com"},
				Success:      true,
			},
			expectedErr: nil,
		},
		{
			name: "NotFound",
			req:  &pb.GetProfessionalRequest{Id: 999},
			mockSetup: func() {
				(mockRepo).On("GetProfessionalByID", uint(999)).Return((*models.Professional)(nil), errors.New("not found")).Once()
			},
			expectedResp: &pb.GetProfessionalResponse{Success: false},
			expectedErr:  errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.GetProfessional(tt.req)
			assert.Equal(t, tt.expectedResp, resp)
			assert.Equal(t, tt.expectedErr, err)
			(mockRepo).AssertExpectations(t)
		})
	}
}

func TestListProfessionalsService(t *testing.T) {
	mockRepo := new(MockProfessionalRepository)
	srv := services.NewProfessionalService(mockRepo)

	tests := []struct {
		name         string
		req          *pb.ListProfessionalsRequest
		mockSetup    func()
		expectedResp *pb.ListProfessionalsResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.ListProfessionalsRequest{},
			mockSetup: func() {
				(mockRepo).On("ListProfessionals").Return([]models.Professional{
					{ID: 1, Name: "Dr. Lopez", Profession: "Dentista", Contact: "lopez@email.com"},
					{ID: 2, Name: "Dr. Perez", Profession: "Medico", Contact: "perez@email.com"},
				}, nil).Once()
			},
			expectedResp: &pb.ListProfessionalsResponse{
				Professionals: []*pb.Professional{
					{Id: 1, Name: "Dr. Lopez", Profession: "Dentista", Contact: "lopez@email.com"},
					{Id: 2, Name: "Dr. Perez", Profession: "Medico", Contact: "perez@email.com"},
				},
				Success: true,
			},
			expectedErr: nil,
		},
		{
			name: "EmptyList",
			req:  &pb.ListProfessionalsRequest{},
			mockSetup: func() {
				(mockRepo).On("ListProfessionals").Return([]models.Professional{}, nil).Once()
			},
			expectedResp: &pb.ListProfessionalsResponse{Professionals: []*pb.Professional{}, Success: true},
			expectedErr:  nil,
		},
		{
			name: "DatabaseError",
			req:  &pb.ListProfessionalsRequest{},
			mockSetup: func() {
				(mockRepo).On("ListProfessionals").Return([]models.Professional(nil), errors.New("db error")).Once()
			},
			expectedResp: &pb.ListProfessionalsResponse{Success: false},
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.ListProfessionals(tt.req)
			assert.Equal(t, tt.expectedResp, resp)
			assert.Equal(t, tt.expectedErr, err)
			(mockRepo).AssertExpectations(t)
		})
	}
}
