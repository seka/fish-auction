package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

func TestAdminRepository_FindOneByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAdminRepository(db)
	email := "admin@example.com"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at"}).
			AddRow(1, email, "hash", time.Now())

		mock.ExpectQuery("SELECT id, email, password_hash, created_at FROM admins WHERE email = \\$1").
			WithArgs(email).
			WillReturnRows(rows)

		got, err := repo.FindOneByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Equal(t, email, got.Email)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, password_hash, created_at FROM admins WHERE email = \\$1").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		got, err := repo.FindOneByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestAdminRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAdminRepository(db)
	id := 1

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at"}).
			AddRow(id, "admin@example.com", "hash", time.Now())

		mock.ExpectQuery("SELECT id, email, password_hash, created_at FROM admins WHERE id = \\$1").
			WithArgs(id).
			WillReturnRows(rows)

		got, err := repo.FindByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, id, got.ID)
	})
}

func TestAdminRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAdminRepository(db)
	admin := &entity.Admin{Email: "new@example.com", PasswordHash: "hash"}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now())

		mock.ExpectQuery("INSERT INTO admins").
			WithArgs(admin.Email, admin.PasswordHash).
			WillReturnRows(rows)

		err := repo.Create(context.Background(), admin)
		assert.NoError(t, err)
		assert.Equal(t, 1, admin.ID)
	})
}

func TestAdminRepository_UpdatePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAdminRepository(db)
	id := 1
	newHash := "newHash"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE admins SET password_hash = \\$1 WHERE id = \\$2").
			WithArgs(newHash, id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdatePassword(context.Background(), id, newHash)
		assert.NoError(t, err)
	})
}
