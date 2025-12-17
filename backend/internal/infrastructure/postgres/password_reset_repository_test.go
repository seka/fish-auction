package postgres_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

func TestPasswordResetRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPasswordResetRepository(db)
	userID := 1
	role := "buyer"
	tokenHash := "hash"
	expiresAt := time.Now()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO password_reset_tokens (user_id, user_role, token_hash, expires_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(userID, role, tokenHash, expiresAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(context.Background(), userID, role, tokenHash, expiresAt)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO password_reset_tokens")).
			WillReturnError(sql.ErrConnDone)

		err := repo.Create(context.Background(), userID, role, tokenHash, expiresAt)
		assert.Error(t, err)
	})
}

func TestPasswordResetRepository_FindByTokenHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPasswordResetRepository(db)
	tokenHash := "hash"
	userID := 1
	role := "buyer"
	expiresAt := time.Now()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"user_id", "user_role", "expires_at"}).
			AddRow(userID, role, expiresAt)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, user_role, expires_at FROM password_reset_tokens WHERE token_hash = $1")).
			WithArgs(tokenHash).
			WillReturnRows(rows)

		gotUserID, gotRole, gotExpiresAt, err := repo.FindByTokenHash(context.Background(), tokenHash)
		assert.NoError(t, err)
		assert.Equal(t, userID, gotUserID)
		assert.Equal(t, role, gotRole)
		// Compare times loosely or exactly? time.Now() serialization might lose precision.
		// Usually Assert.Equal handles time comparison well if zone matches. SQL driver might return slightly different pointer.
		// Let's assume standard behavior. If flaky, check assert.WithinDuration.
		assert.Equal(t, expiresAt, gotExpiresAt)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, user_role, expires_at FROM password_reset_tokens WHERE token_hash = $1")).
			WithArgs(tokenHash).
			WillReturnError(sql.ErrNoRows)

		gotUserID, gotRole, _, err := repo.FindByTokenHash(context.Background(), tokenHash)
		assert.NoError(t, err) // Should return nil error and zero values
		assert.Equal(t, 0, gotUserID)
		assert.Equal(t, "", gotRole)
	})
}

func TestPasswordResetRepository_DeleteByTokenHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPasswordResetRepository(db)
	tokenHash := "hash"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM password_reset_tokens WHERE token_hash = $1")).
			WithArgs(tokenHash).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteByTokenHash(context.Background(), tokenHash)
		assert.NoError(t, err)
	})
}

func TestPasswordResetRepository_DeleteAllByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPasswordResetRepository(db)
	userID := 1
	role := "buyer"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM password_reset_tokens WHERE user_id = $1 AND user_role = $2")).
			WithArgs(userID, role).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteAllByUserID(context.Background(), userID, role)
		assert.NoError(t, err)
	})
}
