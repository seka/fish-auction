package postgres_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/stretchr/testify/assert"
)

func TestFishermanStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewFishermanStore(postgres.NewClient(db))
	name := "Fisherman A"

	mock.ExpectQuery("INSERT INTO fishermen").
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, name))

	created, err := repo.Create(context.Background(), name)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.Equal(t, name, created.Name)
}

func TestFishermanStore_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewFishermanStore(postgres.NewClient(db))

	mock.ExpectQuery("SELECT id, name FROM fishermen").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Fisherman A").
			AddRow(2, "Fisherman B"))

	list, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestFishermanStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewFishermanStore(postgres.NewClient(db))
	id := 1

	mock.ExpectExec("UPDATE fishermen SET deleted_at = CURRENT_TIMESTAMP WHERE id = \\$1").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), id)
	assert.NoError(t, err)
}
