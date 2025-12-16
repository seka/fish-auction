package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAuthenticationRepository(db)
	auth := &model.Authentication{
		BuyerID:      1,
		Email:        "buyer@example.com",
		PasswordHash: "hash",
		AuthType:     "password",
	}

	mock.ExpectQuery("INSERT INTO authentications").
		WithArgs(auth.BuyerID, auth.Email, auth.PasswordHash, auth.AuthType).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	created, err := repo.Create(context.Background(), auth)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
}

func TestAuthenticationRepository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAuthenticationRepository(db)
	email := "buyer@example.com"

	mock.ExpectQuery("SELECT .* FROM authentications WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "buyer_id", "email", "password_hash", "auth_type", "failed_attempts", "locked_until", "last_login_at", "created_at", "updated_at"}).
			AddRow(1, 1, email, "hash", "password", 0, nil, nil, time.Now(), time.Now()))

	found, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, email, found.Email)
}

func TestAuthenticationRepository_UpdatePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAuthenticationRepository(db)
	buyerID := 1
	newHash := "newHash"

	mock.ExpectExec("UPDATE authentications SET password_hash = \\$1, updated_at = CURRENT_TIMESTAMP WHERE buyer_id = \\$2").
		WithArgs(newHash, buyerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdatePassword(context.Background(), buyerID, newHash)
	assert.NoError(t, err)
}
