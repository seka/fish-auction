package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

type mockFishermanCache struct{}

func (m *mockFishermanCache) Get(ctx context.Context, id int) (*model.Fisherman, error) {
	return nil, errors.New("cache miss")
}
func (m *mockFishermanCache) Set(ctx context.Context, id int, fisherman *model.Fisherman) error {
	return nil
}
func (m *mockFishermanCache) Delete(ctx context.Context, id int) error { return nil }

func TestFishermanRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewFishermanRepository(db, &mockFishermanCache{})
	name := "Fisherman A"

	mock.ExpectQuery("INSERT INTO fishermen").
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, name))

	created, err := repo.Create(context.Background(), name)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.Equal(t, name, created.Name)
}

func TestFishermanRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewFishermanRepository(db, &mockFishermanCache{})

	mock.ExpectQuery("SELECT id, name FROM fishermen").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Fisherman A").
			AddRow(2, "Fisherman B"))

	list, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}
