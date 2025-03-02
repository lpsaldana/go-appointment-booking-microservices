package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAgendaRepository struct {
	mock.Mock
}

func (m *MockAgendaRepository) CreateSlot(slot *models.Slot) error {
	args := m.Called(slot)
	return args.Error(0)
}

func (m *MockAgendaRepository) ListAvailableSlots(professionalID uint, date time.Time) ([]models.Slot, error) {
	args := m.Called(professionalID, date)
	return args.Get(0).([]models.Slot), args.Error(1)
}

func (m *MockAgendaRepository) CreateAppointment(appointment *models.Appointment) error {
	args := m.Called(appointment)
	return args.Error(0)
}

func (m *MockAgendaRepository) UpdateSlotAvailability(slotID uint, available bool) error {
	args := m.Called(slotID, available)
	return args.Error(0)
}

func (m *MockAgendaRepository) ListAppointments(clientID, professionalID uint) ([]models.Appointment, error) {
	args := m.Called(clientID, professionalID)
	return args.Get(0).([]models.Appointment), args.Error(1)
}

func (m *MockAgendaRepository) GetSlotByID(slotID uint) (*models.Slot, error) {
	args := m.Called(slotID)
	return args.Get(0).(*models.Slot), args.Error(1)
}

// Mock para NotificationServiceClient
type MockNotificationServiceClient struct {
	mock.Mock
}

func (m *MockNotificationServiceClient) SendAppointmentNotification(ctx context.Context, in *pb.SendAppointmentNotificationRequest, opts ...grpc.CallOption) (*pb.SendAppointmentNotificationResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.SendAppointmentNotificationResponse), args.Error(1)
}

func TestCreateSlot(t *testing.T) {
	mockRepo := new(MockAgendaRepository)
	mockNotif := new(MockNotificationServiceClient)
	// Creamos el servicio con un *grpc.ClientConn dummy (nil), y luego inyectamos el mock
	srv := services.NewAgendaService(mockRepo, nil)
	srv.(*services.AgendaServiceImpl).NotifClient = mockNotif // Inyectamos el mock después

	tests := []struct {
		name         string
		req          *pb.CreateSlotRequest
		mockSetup    func()
		expectedResp *pb.CreateSlotResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.CreateSlotRequest{ProfessionalId: 1, StartTime: "2025-03-10T10:00:00Z", EndTime: "2025-03-10T10:30:00Z"},
			mockSetup: func() {
				(mockRepo).On("CreateSlot", mock.AnythingOfType("*models.Slot")).Return(nil).Once()
			},
			expectedResp: &pb.CreateSlotResponse{Message: "Slot created", Success: true, SlotId: 0},
			expectedErr:  nil,
		},
		{
			name:         "InvalidStartTime",
			req:          &pb.CreateSlotRequest{ProfessionalId: 1, StartTime: "invalid", EndTime: "2025-03-10T10:30:00Z"},
			mockSetup:    func() {},
			expectedResp: &pb.CreateSlotResponse{Message: "start_time invalid format", Success: false},
			expectedErr:  errors.New("parsing time \"invalid\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid\" as \"2006\""),
		},
		{
			name: "DatabaseError",
			req:  &pb.CreateSlotRequest{ProfessionalId: 1, StartTime: "2025-03-10T10:00:00Z", EndTime: "2025-03-10T10:30:00Z"},
			mockSetup: func() {
				(mockRepo).On("CreateSlot", mock.AnythingOfType("*models.Slot")).Return(errors.New("db error")).Once()
			},
			expectedResp: &pb.CreateSlotResponse{Message: "Error creating slot", Success: false},
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.CreateSlot(tt.req)
			assert.Equal(t, tt.expectedResp.Message, resp.Message)
			assert.Equal(t, tt.expectedResp.Success, resp.Success)
			/*if tt.expectedResp.Success {
				assert.NotZero(t, resp.SlotId, "SlotId debería asignarse")
			}*/
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			(mockRepo).AssertExpectations(t)
			(mockNotif).AssertExpectations(t)
		})
	}
}

func TestListAvailableSlots(t *testing.T) {
	mockRepo := new(MockAgendaRepository)
	mockNotif := new(MockNotificationServiceClient)
	srv := services.NewAgendaService(mockRepo, nil)
	srv.(*services.AgendaServiceImpl).NotifClient = mockNotif

	tests := []struct {
		name         string
		req          *pb.ListAvailableSlotsRequest
		mockSetup    func()
		expectedResp *pb.ListAvailableSlotsResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.ListAvailableSlotsRequest{ProfessionalId: 1, Date: "2025-03-10"},
			mockSetup: func() {
				start, _ := time.Parse("2006-01-02", "2025-03-10")
				(mockRepo).On("ListAvailableSlots", uint(1), start).Return([]models.Slot{
					{ID: 1, ProfessionalID: 1, StartTime: time.Now(), EndTime: time.Now().Add(30 * time.Minute), Available: true},
				}, nil).Once()
			},
			expectedResp: &pb.ListAvailableSlotsResponse{
				Slots: []*pb.Slot{
					{Id: 1, ProfessionalId: 1, StartTime: time.Now().Format(time.RFC3339), EndTime: time.Now().Add(30 * time.Minute).Format(time.RFC3339), Available: true},
				},
				Success: true,
			},
			expectedErr: nil,
		},
		{
			name:         "InvalidDate",
			req:          &pb.ListAvailableSlotsRequest{ProfessionalId: 1, Date: "invalid"},
			mockSetup:    func() {},
			expectedResp: &pb.ListAvailableSlotsResponse{Success: false},
			expectedErr:  errors.New("parsing time \"invalid\" as \"2006-01-02\": cannot parse \"invalid\" as \"2006\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.ListAvailableSlots(tt.req)
			if tt.expectedResp.Success {
				assert.Equal(t, tt.expectedResp.Success, resp.Success)
				assert.Len(t, resp.Slots, len(tt.expectedResp.Slots))
				for i, slot := range resp.Slots {
					assert.Equal(t, tt.expectedResp.Slots[i].Id, slot.Id)
					assert.Equal(t, tt.expectedResp.Slots[i].ProfessionalId, slot.ProfessionalId)
					assert.Equal(t, tt.expectedResp.Slots[i].Available, slot.Available)
				}
			} else {
				assert.Equal(t, tt.expectedResp, resp)
			}
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			(mockRepo).AssertExpectations(t)
			(mockNotif).AssertExpectations(t)
		})
	}
}

func TestBookAppointment(t *testing.T) {
	mockRepo := new(MockAgendaRepository)
	mockNotif := new(MockNotificationServiceClient)
	srv := services.NewAgendaService(mockRepo, nil)
	srv.(*services.AgendaServiceImpl).NotifClient = mockNotif // Inyectamos el mock después

	tests := []struct {
		name         string
		req          *pb.BookAppointmentRequest
		mockSetup    func()
		expectedResp *pb.BookAppointmentResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.BookAppointmentRequest{ClientId: 1, SlotId: 1},
			mockSetup: func() {
				(mockRepo).On("GetSlotByID", uint(1)).Return(&models.Slot{ID: 1, ProfessionalID: 2, Available: true}, nil).Once()
				(mockRepo).On("CreateAppointment", mock.AnythingOfType("*models.Appointment")).Return(nil).Once()
				(mockRepo).On("UpdateSlotAvailability", uint(1), false).Return(nil).Once()
				(mockNotif).On("SendAppointmentNotification", nil, mock.AnythingOfType("*pb.SendAppointmentNotificationRequest")).
					Return(&pb.SendAppointmentNotificationResponse{Message: "Sent", Success: true}, nil).Once()
			},
			expectedResp: &pb.BookAppointmentResponse{Message: "Cita reservada exitosamente", Success: true, AppointmentId: 1},
			expectedErr:  nil,
		},
		{
			name: "SlotNotFound",
			req:  &pb.BookAppointmentRequest{ClientId: 1, SlotId: 999},
			mockSetup: func() {
				(mockRepo).On("GetSlotByID", uint(999)).Return((*models.Slot)(nil), errors.New("not found")).Once()
			},
			expectedResp: &pb.BookAppointmentResponse{Message: "Slot no encontrado", Success: false},
			expectedErr:  errors.New("not found"),
		},
		{
			name: "SlotNotAvailable",
			req:  &pb.BookAppointmentRequest{ClientId: 1, SlotId: 1},
			mockSetup: func() {
				(mockRepo).On("GetSlotByID", uint(1)).Return(&models.Slot{ID: 1, ProfessionalID: 2, Available: false}, nil).Once()
			},
			expectedResp: &pb.BookAppointmentResponse{Message: "El slot no está disponible", Success: false},
			expectedErr:  nil,
		},
		{
			name: "NotificationError",
			req:  &pb.BookAppointmentRequest{ClientId: 1, SlotId: 1},
			mockSetup: func() {
				(mockRepo).On("GetSlotByID", uint(1)).Return(&models.Slot{ID: 1, ProfessionalID: 2, Available: true}, nil).Once()
				(mockRepo).On("CreateAppointment", mock.AnythingOfType("*models.Appointment")).Return(nil).Once()
				(mockRepo).On("UpdateSlotAvailability", uint(1), false).Return(nil).Once()
				(mockNotif).On("SendAppointmentNotification", nil, mock.AnythingOfType("*pb.SendAppointmentNotificationRequest")).
					Return(&pb.SendAppointmentNotificationResponse{Message: "Error", Success: false}, errors.New("notification failed")).Once()
			},
			expectedResp: &pb.BookAppointmentResponse{Message: "Cita reservada exitosamente", Success: true, AppointmentId: 1},
			expectedErr:  nil, // Error de notificación no afecta la reserva
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.BookAppointment(tt.req)
			assert.Equal(t, tt.expectedResp.Message, resp.Message)
			assert.Equal(t, tt.expectedResp.Success, resp.Success)
			if tt.expectedResp.Success {
				assert.NotZero(t, resp.AppointmentId, "AppointmentId debería asignarse")
			}
			assert.Equal(t, tt.expectedErr, err)
			(mockRepo).AssertExpectations(t)
			(mockNotif).AssertExpectations(t)
		})
	}
}

func TestListAppointments(t *testing.T) {
	mockRepo := new(MockAgendaRepository)
	mockNotif := new(MockNotificationServiceClient)
	srv := services.NewAgendaService(mockRepo, nil)
	srv.(*services.AgendaServiceImpl).NotifClient = mockNotif

	tests := []struct {
		name         string
		req          *pb.ListAppointmentsRequest
		mockSetup    func()
		expectedResp *pb.ListAppointmentsResponse
		expectedErr  error
	}{
		{
			name: "Success",
			req:  &pb.ListAppointmentsRequest{ClientId: 1},
			mockSetup: func() {
				(mockRepo).On("ListAppointments", uint(1), uint(0)).Return([]models.Appointment{
					{ID: 1, ClientID: 1, SlotID: 1, ProfessionalID: 2},
				}, nil).Once()
				(mockRepo).On("GetSlotByID", uint(1)).Return(&models.Slot{ID: 1, StartTime: time.Now(), EndTime: time.Now().Add(30 * time.Minute)}, nil).Once()
			},
			expectedResp: &pb.ListAppointmentsResponse{
				Appointments: []*pb.Appointment{
					{Id: 1, ClientId: 1, SlotId: 1, ProfessionalId: 2, StartTime: time.Now().Format(time.RFC3339), EndTime: time.Now().Add(30 * time.Minute).Format(time.RFC3339)},
				},
				Success: true,
			},
			expectedErr: nil,
		},
		{
			name: "EmptyList",
			req:  &pb.ListAppointmentsRequest{ClientId: 1},
			mockSetup: func() {
				(mockRepo).On("ListAppointments", uint(1), uint(0)).Return([]models.Appointment{}, nil).Once()
			},
			expectedResp: &pb.ListAppointmentsResponse{Appointments: []*pb.Appointment{}, Success: true},
			expectedErr:  nil,
		},
		{
			name: "DatabaseError",
			req:  &pb.ListAppointmentsRequest{ClientId: 1},
			mockSetup: func() {
				(mockRepo).On("ListAppointments", uint(1), uint(0)).Return([]models.Appointment{}, errors.New("db error")).Once()
			},
			expectedResp: &pb.ListAppointmentsResponse{Success: false},
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := srv.ListAppointments(tt.req)
			assert.Equal(t, tt.expectedResp.Success, resp.Success)
			assert.Len(t, resp.Appointments, len(tt.expectedResp.Appointments))
			for i, appt := range resp.Appointments {
				assert.Equal(t, tt.expectedResp.Appointments[i].Id, appt.Id)
				assert.Equal(t, tt.expectedResp.Appointments[i].ClientId, appt.ClientId)
				assert.Equal(t, tt.expectedResp.Appointments[i].ProfessionalId, appt.ProfessionalId)
			}
			assert.Equal(t, tt.expectedErr, err)
			(mockRepo).AssertExpectations(t)
			(mockNotif).AssertExpectations(t)
		})
	}
}
