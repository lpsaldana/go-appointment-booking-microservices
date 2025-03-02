package unit

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repositories.UserRepository) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	assert.NoError(t, err)
	repo := repositories.NewUserRepository(gormDB)
	return sqlDB, mock, repo
}

func TestCreateUserRepo(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name      string
		user      *models.User // Movemos la definición del user aquí
		mockSetup func(sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "Success",
			user: &models.User{Username: "testuser", Password: "testpass"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("username","password") VALUES ($1,$2) RETURNING "id"`)).
					WithArgs("testuser", sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "DatabaseError",
			user: &models.User{Username: "testuser", Password: "testpass"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("username","password") VALUES ($1,$2) RETURNING "id"`)).
					WithArgs("testuser", sqlmock.AnyArg()).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			err := repo.CreateUser(tt.user)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, "testpass", tt.user.Password, "La contraseña debería haberse encriptado")
				assert.Contains(t, tt.user.Password, "$2a$", "La contraseña debería ser un hash bcrypt")
				assert.Equal(t, uint(1), tt.user.ID, "El ID debería haberse asignado")
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFindByUsername(t *testing.T) {
	sqlDB, mock, repo := setupMockDB(t)
	defer sqlDB.Close()

	tests := []struct {
		name         string
		username     string
		mockSetup    func(sqlmock.Sqlmock)
		expectedUser *models.User
		expectedErr  error
	}{
		{
			name:     "Success",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(1, "testuser", "hashedpass")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("testuser", 1).
					WillReturnRows(rows)
			},
			expectedUser: &models.User{ID: 1, Username: "testuser", Password: "hashedpass"},
			expectedErr:  nil,
		},
		{
			name:     "NotFound",
			username: "unknown",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("unknown", 1).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser: nil,
			expectedErr:  sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			user, err := repo.FindByUsername(tt.username)
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
