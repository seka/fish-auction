package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
)

func TestAdminPasswordResetRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAdminPasswordResetRepository(db)

	adminID := 1
	tokenHash := "hash123"
	expiresAt := time.Now()

	mock.ExpectExec("INSERT INTO admin_password_reset_tokens").
		WithArgs(adminID, tokenHash, expiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), adminID, tokenHash, expiresAt)
	if err != nil {
		t.Errorf("error was not expected while updates: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAdminPasswordResetRepository_FindByTokenHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAdminPasswordResetRepository(db)

	tokenHash := "hash123"
	adminID := 1
	expiresAt := time.Now()

	// Case 1: Found
	rows := sqlmock.NewRows([]string{"admin_id", "expires_at"}).
		AddRow(adminID, expiresAt)

	mock.ExpectQuery("SELECT admin_id, expires_at FROM admin_password_reset_tokens").
		WithArgs(tokenHash).
		WillReturnRows(rows)

	gotID, gotExpires, err := repo.FindByTokenHash(context.Background(), tokenHash)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if gotID != adminID {
		t.Errorf("expected adminID %d, got %d", adminID, gotID)
	}
	if !gotExpires.Equal(expiresAt) {
		t.Errorf("expected expiresAt %v, got %v", expiresAt, gotExpires)
	}

	// Case 2: Not Found
	mock.ExpectQuery("SELECT admin_id, expires_at FROM admin_password_reset_tokens").
		WithArgs("invalid").
		WillReturnError(sql.ErrNoRows)

	gotID, _, err = repo.FindByTokenHash(context.Background(), "invalid")
	if err != nil {
		t.Errorf("expected no error on not found, got %v", err)
	}
	if gotID != 0 {
		t.Errorf("expected 0 ID, got %d", gotID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAdminPasswordResetRepository_DeleteByTokenHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAdminPasswordResetRepository(db)
	tokenHash := "hash123"

	mock.ExpectExec("DELETE FROM admin_password_reset_tokens").
		WithArgs(tokenHash).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteByTokenHash(context.Background(), tokenHash)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
