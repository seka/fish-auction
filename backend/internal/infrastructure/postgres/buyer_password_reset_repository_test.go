package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
)

func TestBuyerPasswordResetRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewBuyerPasswordResetRepository(db)

	buyerID := 1
	tokenHash := "hash123"
	expiresAt := time.Now()

	mock.ExpectExec("INSERT INTO buyer_password_reset_tokens").
		WithArgs(buyerID, tokenHash, expiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), buyerID, tokenHash, expiresAt)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBuyerPasswordResetRepository_FindByTokenHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewBuyerPasswordResetRepository(db)

	tokenHash := "hash123"
	buyerID := 1
	expiresAt := time.Now()

	// Case 1: Found
	rows := sqlmock.NewRows([]string{"buyer_id", "expires_at"}).
		AddRow(buyerID, expiresAt)

	mock.ExpectQuery("SELECT buyer_id, expires_at FROM buyer_password_reset_tokens").
		WithArgs(tokenHash).
		WillReturnRows(rows)

	gotID, gotExpires, err := repo.FindByTokenHash(context.Background(), tokenHash)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if gotID != buyerID {
		t.Errorf("expected buyerID %d, got %d", buyerID, gotID)
	}
	if !gotExpires.Equal(expiresAt) {
		t.Errorf("expected expiresAt %v, got %v", expiresAt, gotExpires)
	}

	// Case 2: Not Found (Repo returns zero values nil error on ErrNoRows)
	mock.ExpectQuery("SELECT buyer_id, expires_at FROM buyer_password_reset_tokens").
		WithArgs("invalid").
		WillReturnError(sql.ErrNoRows)

	gotID, _, err = repo.FindByTokenHash(context.Background(), "invalid")
	if err != nil {
		t.Errorf("expected no error on not found, got %v", err)
	}
	if gotID != 0 {
		t.Errorf("expected 0 ID, got %d", gotID)
	}

	// Case 3: Error
	mock.ExpectQuery("SELECT buyer_id, expires_at FROM buyer_password_reset_tokens").
		WithArgs("error_case").
		WillReturnError(errors.New("db error"))

	_, _, err = repo.FindByTokenHash(context.Background(), "error_case")
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBuyerPasswordResetRepository_DeleteByTokenHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewBuyerPasswordResetRepository(db)
	tokenHash := "hash123"

	mock.ExpectExec("DELETE FROM buyer_password_reset_tokens").
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
