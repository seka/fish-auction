package postgres_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/stretchr/testify/assert"
)

func TestPasswordResetStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewPasswordResetStore(postgres.NewClient(db))
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

func TestPasswordResetStore_FindByTokenHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewPasswordResetStore(postgres.NewClient(db))
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

		got, err := repo.FindByTokenHash(context.Background(), tokenHash)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, userID, got.UserID)
		assert.Equal(t, role, got.Role)
		assert.Equal(t, tokenHash, got.TokenHash)
		assert.Equal(t, expiresAt.Unix(), got.ExpiresAt.Unix()) // Compare unix timestamps to avoid precision issues
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, user_role, expires_at FROM password_reset_tokens WHERE token_hash = $1")).
			WithArgs(tokenHash).
			WillReturnError(sql.ErrNoRows)

		got, err := repo.FindByTokenHash(context.Background(), tokenHash)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestPasswordResetStore_DeleteByTokenHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewPasswordResetStore(postgres.NewClient(db))
	tokenHash := "hash"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM password_reset_tokens WHERE token_hash = $1")).
			WithArgs(tokenHash).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteByTokenHash(context.Background(), tokenHash)
		assert.NoError(t, err)
	})
}

func TestPasswordResetStore_DeleteAllByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewPasswordResetStore(postgres.NewClient(db))
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
