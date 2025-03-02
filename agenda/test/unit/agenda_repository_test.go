package unit

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repositories.AgendaRepository) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	assert.NoError(t, err)
	repo := repositories.NewAgendaRepository(gormDB)
	return sqlDB, mock, repo
}

func TestCreateSlotRepo(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name      string
		slot      *models.Slot
		mockSetup func(sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "Success",
			slot: &models.Slot{ProfessionalID: 1, StartTime: time.Now(), EndTime: time.Now().Add(30 * time.Minute), Available: true},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "slots" ("professional_id","start_time","end_time","available") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(uint(1), sqlmock.AnyArg(), sqlmock.AnyArg(), true).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "DatabaseError",
			slot: &models.Slot{ProfessionalID: 1, StartTime: time.Now(), EndTime: time.Now().Add(30 * time.Minute), Available: true},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "slots" ("professional_id","start_time","end_time","available") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(uint(1), sqlmock.AnyArg(), sqlmock.AnyArg(), true).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			err := repo.CreateSlot(tt.slot)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint(1), tt.slot.ID, "El ID debería haberse asignado")
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListAvailableSlotsRepo(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name           string
		professionalID uint
		date           time.Time
		mockSetup      func(sqlmock.Sqlmock)
		expectedSlots  []models.Slot
		expectedErr    error
	}{
		{
			name:           "Success",
			professionalID: 1,
			date:           time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Parseamos las fechas como time.Time para que coincidan con el modelo
				startTime, _ := time.Parse(time.RFC3339, "2025-03-10T10:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-03-10T10:30:00Z")
				rows := sqlmock.NewRows([]string{"id", "professional_id", "start_time", "end_time", "available"}).
					AddRow(1, 1, startTime, endTime, true)
				startOfDay := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)
				endOfDay := startOfDay.Add(24 * time.Hour)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "slots" WHERE professional_id = $1 AND start_time >= $2 AND start_time < $3 AND available = $4`)).
					WithArgs(1, startOfDay, endOfDay, true).
					WillReturnRows(rows)
			},
			expectedSlots: []models.Slot{
				{ID: 1, ProfessionalID: 1, StartTime: time.Date(2025, 3, 10, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 3, 10, 10, 30, 0, 0, time.UTC), Available: true},
			},
			expectedErr: nil,
		},
		{
			name:           "EmptyList",
			professionalID: 1,
			date:           time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "professional_id", "start_time", "end_time", "available"})
				startOfDay := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)
				endOfDay := startOfDay.Add(24 * time.Hour)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "slots" WHERE professional_id = $1 AND start_time >= $2 AND start_time < $3 AND available = $4`)).
					WithArgs(1, startOfDay, endOfDay, true).
					WillReturnRows(rows)
			},
			expectedSlots: []models.Slot{},
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			slots, err := repo.ListAvailableSlots(tt.professionalID, tt.date)
			assert.Equal(t, tt.expectedSlots, slots)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateAppointment(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name        string
		appointment *models.Appointment
		mockSetup   func(sqlmock.Sqlmock)
		expectErr   bool
	}{
		{
			name:        "Success",
			appointment: &models.Appointment{ClientID: 1, SlotID: 1, ProfessionalID: 2},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "appointments" ("client_id","slot_id","professional_id") VALUES ($1,$2,$3) RETURNING "id"`)).
					WithArgs(uint(1), uint(1), uint(2)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:        "DatabaseError",
			appointment: &models.Appointment{ClientID: 1, SlotID: 1, ProfessionalID: 2},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "appointments" ("client_id","slot_id","professional_id") VALUES ($1,$2,$3) RETURNING "id"`)).
					WithArgs(uint(1), uint(1), uint(2)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			err := repo.CreateAppointment(tt.appointment)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint(1), tt.appointment.ID, "El ID debería haberse asignado")
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateSlotAvailability(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name      string
		slotID    uint
		available bool
		mockSetup func(sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:      "Success",
			slotID:    1,
			available: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "slots" SET "available"=$1 WHERE id = $2`)).
					WithArgs(false, uint(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:      "DatabaseError",
			slotID:    1,
			available: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "slots" SET "available"=$1 WHERE id = $2`)).
					WithArgs(false, uint(1)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			err := repo.UpdateSlotAvailability(tt.slotID, tt.available)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListAppointmentsRepo(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name           string
		clientID       uint
		professionalID uint
		mockSetup      func(sqlmock.Sqlmock)
		expectedAppts  []models.Appointment
		expectedErr    error
	}{
		{
			name:           "SuccessWithClientID",
			clientID:       1,
			professionalID: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "slot_id", "professional_id"}).
					AddRow(1, 1, 1, 2)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "appointments" WHERE client_id = $1`)).
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			expectedAppts: []models.Appointment{{ID: 1, ClientID: 1, SlotID: 1, ProfessionalID: 2}},
			expectedErr:   nil,
		},
		{
			name:           "EmptyList",
			clientID:       1,
			professionalID: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "slot_id", "professional_id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "appointments" WHERE client_id = $1`)).
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			expectedAppts: []models.Appointment{},
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			appointments, err := repo.ListAppointments(tt.clientID, tt.professionalID)
			assert.Equal(t, tt.expectedAppts, appointments)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetSlotByID(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name         string
		slotID       uint
		mockSetup    func(sqlmock.Sqlmock)
		expectedSlot *models.Slot
		expectedErr  error
	}{
		{
			name:   "Success",
			slotID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				startTime, _ := time.Parse(time.RFC3339, "2025-03-10T10:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-03-10T10:30:00Z")
				rows := sqlmock.NewRows([]string{"id", "professional_id", "start_time", "end_time", "available"}).
					AddRow(1, 2, startTime, endTime, true)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "slots" WHERE "slots"."id" = $1 ORDER BY "slots"."id" LIMIT $2`)).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnRows(rows)
			},
			expectedSlot: &models.Slot{ID: 1, ProfessionalID: 2, StartTime: time.Date(2025, 3, 10, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 3, 10, 10, 30, 0, 0, time.UTC), Available: true},
			expectedErr:  nil,
		},
		{
			name:   "NotFound",
			slotID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "slots" WHERE "slots"."id" = $1 ORDER BY "slots"."id" LIMIT $2`)).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(sql.ErrNoRows)
			},
			expectedSlot: nil,
			expectedErr:  sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			slot, err := repo.GetSlotByID(tt.slotID)
			assert.Equal(t, tt.expectedSlot, slot)
			assert.Equal(t, tt.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
