package unit

import (
	"errors"
	"testing"

	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) CreateClient(client *models.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientRepository) GetClientByID(id uint) (*models.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientRepository) ListClients() ([]models.Client, error) {
	args := m.Called()
	return args.Get(0).([]models.Client), args.Error(1)
}

func TestCreateClient(t *testing.T) {
	mockRepo := new(MockClientRepository)
	srv := services.NewClientService(mockRepo)

	tests := []struct {
		name         string
		req          *pb.CreateClientRequest
		mockSetup    func()
		expectedResp *pb.CreateClientResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.CreateClientRequest{Name: "Maria Perez", Email: "maria@email.com", Phone: "123456789"},
			mockSetup: func() {
				(mockRepo).On("CreateClient", mock.AnythingOfType("*models.Client")).Return(nil).Once()
			},
			expectedResp: &pb.CreateClientResponse{Message: "Client created", Success: true, ClientId: 0},
			expectedErr:  nil,
		},
		{
			name: "DatabaseError",
			req:  &pb.CreateClientRequest{Name: "Pedro Gomez", Email: "pedro@email.com", Phone: "987654321"},
			mockSetup: func() {
				(mockRepo).On("CreateClient", mock.AnythingOfType("*models.Client")).Return(errors.New("db error")).Once()
			},
			expectedResp: &pb.CreateClientResponse{Message: "Error creating client", Success: false},
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.CreateClient(tt.req)
			assert.Equal(t, tt.expectedResp.Message, resp.Message)
			assert.Equal(t, tt.expectedResp.Success, resp.Success)
			/*if tt.expectedResp.Success {
				assert.NotZero(t, resp.ClientId, "ClientId debería asignarse")
			}*/
			assert.Equal(t, tt.expectedErr, err)
			(mockRepo).AssertExpectations(t)
		})
	}
}

func TestGetClient(t *testing.T) {
	mockRepo := new(MockClientRepository)
	srv := services.NewClientService(mockRepo)

	tests := []struct {
		name         string
		req          *pb.GetClientRequest
		mockSetup    func()
		expectedResp *pb.GetClientResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.GetClientRequest{Id: 1},
			mockSetup: func() {
				(mockRepo).On("GetClientByID", uint(1)).Return(&models.Client{ID: 1, Name: "Maria Perez", Email: "maria@email.com", Phone: "123456789"}, nil).Once()
			},
			expectedResp: &pb.GetClientResponse{
				Client:  &pb.Client{Id: 1, Name: "Maria Perez", Email: "maria@email.com", Phone: "123456789"},
				Success: true,
			},
			expectedErr: nil,
		},
		{
			name: "NotFound",
			req:  &pb.GetClientRequest{Id: 999},
			mockSetup: func() {
				(mockRepo).On("GetClientByID", uint(999)).Return((*models.Client)(nil), errors.New("not found")).Once()
			},
			expectedResp: &pb.GetClientResponse{Success: false},
			expectedErr:  errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.GetClient(tt.req)
			assert.Equal(t, tt.expectedResp, resp)
			assert.Equal(t, tt.expectedErr, err)
			(mockRepo).AssertExpectations(t)
		})
	}
}

func TestListClients(t *testing.T) {
	mockRepo := new(MockClientRepository)
	srv := services.NewClientService(mockRepo)

	tests := []struct {
		name         string
		req          *pb.ListClientsRequest
		mockSetup    func()
		expectedResp *pb.ListClientsResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.ListClientsRequest{},
			mockSetup: func() {
				(mockRepo).On("ListClients").Return([]models.Client{
					{ID: 1, Name: "Maria Perez", Email: "maria@email.com", Phone: "123456789"},
					{ID: 2, Name: "Pedro Gomez", Email: "pedro@email.com", Phone: "987654321"},
				}, nil).Once()
			},
			expectedResp: &pb.ListClientsResponse{
				Clients: []*pb.Client{
					{Id: 1, Name: "Maria Perez", Email: "maria@email.com", Phone: "123456789"},
					{Id: 2, Name: "Pedro Gomez", Email: "pedro@email.com", Phone: "987654321"},
				},
				Success: true,
			},
			expectedErr: nil,
		},
		{
			name: "EmptyList",
			req:  &pb.ListClientsRequest{},
			mockSetup: func() {
				(mockRepo).On("ListClients").Return([]models.Client{}, nil).Once()
			},
			expectedResp: &pb.ListClientsResponse{Clients: []*pb.Client{}, Success: true},
			expectedErr:  nil,
		},
		{
			name: "DatabaseError",
			req:  &pb.ListClientsRequest{},
			mockSetup: func() {
				(mockRepo).On("ListClients").Return([]models.Client(nil), errors.New("db error")).Once()
			},
			expectedResp: &pb.ListClientsResponse{Success: false},
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.ListClients(tt.req)
			assert.Equal(t, tt.expectedResp, resp)
			assert.Equal(t, tt.expectedErr, err)
			(mockRepo).AssertExpectations(t)
		})
	}
}
