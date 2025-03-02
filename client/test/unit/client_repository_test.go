package unit

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repositories.ClientRepository) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	assert.NoError(t, err)
	repo := repositories.NewClientRepository(gormDB)
	return sqlDB, mock, repo
}

func TestCreateClientRepo(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name      string
		client    *models.Client
		mockSetup func(sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:   "Success",
			client: &models.Client{Name: "Maria Perez", Email: "maria@email.com", Phone: "123456789"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "clients" ("name","email","phone") VALUES ($1,$2,$3) RETURNING "id"`)).
					WithArgs("Maria Perez", "maria@email.com", "123456789").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:   "DatabaseError",
			client: &models.Client{Name: "Pedro Gomez", Email: "pedro@email.com", Phone: "987654321"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "clients" ("name","email","phone") VALUES ($1,$2,$3) RETURNING "id"`)).
					WithArgs("Pedro Gomez", "pedro@email.com", "987654321").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			err := repo.CreateClient(tt.client)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint(1), tt.client.ID, "El ID debería haberse asignado")
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetClientByID(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name           string
		id             uint
		mockSetup      func(sqlmock.Sqlmock)
		expectedClient *models.Client
		expectedErr    error
	}{
		{
			name: "Success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
					AddRow(1, "Maria Perez", "maria@email.com", "123456789")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "clients" WHERE "clients"."id" = $1 ORDER BY "clients"."id" LIMIT $2`)).
					WithArgs(sqlmock.AnyArg(), 1). // Dos argumentos: id y LIMIT
					WillReturnRows(rows)
			},
			expectedClient: &models.Client{ID: 1, Name: "Maria Perez", Email: "maria@email.com", Phone: "123456789"},
			expectedErr:    nil,
		},
		{
			name: "NotFound",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "clients" WHERE "clients"."id" = $1 ORDER BY "clients"."id" LIMIT $2`)).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(sql.ErrNoRows)
			},
			expectedClient: nil,
			expectedErr:    sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			client, err := repo.GetClientByID(tt.id)
			assert.Equal(t, tt.expectedClient, client)
			assert.Equal(t, tt.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListClientsRepo(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name            string
		mockSetup       func(sqlmock.Sqlmock)
		expectedClients []models.Client
		expectedErr     error
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
					AddRow(1, "Maria Perez", "maria@email.com", "123456789").
					AddRow(2, "Pedro Gomez", "pedro@email.com", "987654321")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "clients"`)).
					WillReturnRows(rows)
			},
			expectedClients: []models.Client{
				{ID: 1, Name: "Maria Perez", Email: "maria@email.com", Phone: "123456789"},
				{ID: 2, Name: "Pedro Gomez", Email: "pedro@email.com", Phone: "987654321"},
			},
			expectedErr: nil,
		},
		{
			name: "EmptyList",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "clients"`)).
					WillReturnRows(rows)
			},
			expectedClients: []models.Client{},
			expectedErr:     nil,
		},
		{
			name: "DatabaseError",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "clients"`)).
					WillReturnError(errors.New("db error"))
			},
			expectedClients: nil,
			expectedErr:     errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			clients, err := repo.ListClients()
			assert.Equal(t, tt.expectedClients, clients)
			assert.Equal(t, tt.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
