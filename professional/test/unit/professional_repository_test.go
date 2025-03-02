package unit

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repositories.ProfessionalRepository) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	assert.NoError(t, err)
	repo := repositories.NewProfessionalRepository(gormDB)
	return sqlDB, mock, repo
}

func TestCreateProfessional(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name         string
		professional *models.Professional
		mockSetup    func(sqlmock.Sqlmock)
		expectErr    bool
	}{
		{
			name:         "Success",
			professional: &models.Professional{Name: "Dr. Lopez", Profession: "Dentista", Contact: "lopez@email.com"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "professionals" ("name","profession","contact") VALUES ($1,$2,$3) RETURNING "id"`)).
					WithArgs("Dr. Lopez", "Dentista", "lopez@email.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:         "DatabaseError",
			professional: &models.Professional{Name: "Dr. Perez", Profession: "Medico", Contact: "perez@email.com"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "professionals" ("name","profession","contact") VALUES ($1,$2,$3) RETURNING "id"`)).
					WithArgs("Dr. Perez", "Medico", "perez@email.com").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			err := repo.CreateProfessional(tt.professional)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint(1), tt.professional.ID, "El ID debería haberse asignado")
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetProfessionalByID(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name         string
		id           uint
		mockSetup    func(sqlmock.Sqlmock)
		expectedProf *models.Professional
		expectedErr  error
	}{
		{
			name: "Success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "profession", "contact"}).
					AddRow(1, "Dr. Lopez", "Dentista", "lopez@email.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "professionals" WHERE "professionals"."id" = $1 ORDER BY "professionals"."id" LIMIT $2`)).
					WithArgs(sqlmock.AnyArg(), 1). // Dos argumentos: id y LIMIT
					WillReturnRows(rows)
			},
			expectedProf: &models.Professional{ID: 1, Name: "Dr. Lopez", Profession: "Dentista", Contact: "lopez@email.com"},
			expectedErr:  nil,
		},
		{
			name: "NotFound",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "professionals" WHERE "professionals"."id" = $1 ORDER BY "professionals"."id" LIMIT $2`)).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(sql.ErrNoRows)
			},
			expectedProf: nil,
			expectedErr:  sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			prof, err := repo.GetProfessionalByID(tt.id)
			assert.Equal(t, tt.expectedProf, prof)
			assert.Equal(t, tt.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListProfessionals(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		expectedProfs []models.Professional
		expectedErr   error
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "profession", "contact"}).
					AddRow(1, "Dr. Lopez", "Dentista", "lopez@email.com").
					AddRow(2, "Dr. Perez", "Medico", "perez@email.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "professionals"`)).
					WillReturnRows(rows)
			},
			expectedProfs: []models.Professional{
				{ID: 1, Name: "Dr. Lopez", Profession: "Dentista", Contact: "lopez@email.com"},
				{ID: 2, Name: "Dr. Perez", Profession: "Medico", Contact: "perez@email.com"},
			},
			expectedErr: nil,
		},
		{
			name: "EmptyList",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "profession", "contact"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "professionals"`)).
					WillReturnRows(rows)
			},
			expectedProfs: []models.Professional{},
			expectedErr:   nil,
		},
		{
			name: "DatabaseError",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "professionals"`)).
					WillReturnError(errors.New("db error"))
			},
			expectedProfs: nil,
			expectedErr:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			profs, err := repo.ListProfessionals()
			assert.Equal(t, tt.expectedProfs, profs)
			assert.Equal(t, tt.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
