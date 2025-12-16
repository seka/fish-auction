package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

func TestVenueRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewVenueRepository(db)
	venue := &model.Venue{Name: "Venue A", Location: "Loc A", Description: "Desc A"}

	mock.ExpectQuery("INSERT INTO venues").
		WithArgs(venue.Name, venue.Location, venue.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location", "description", "created_at"}).
			AddRow(1, venue.Name, venue.Location, venue.Description, time.Now()))

	created, err := repo.Create(context.Background(), venue)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
}

func TestVenueRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewVenueRepository(db)

	mock.ExpectQuery("SELECT id, name, location, description, created_at FROM venues ORDER BY created_at DESC").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location", "description", "created_at"}).
			AddRow(1, "Venue A", "Loc A", "Desc A", time.Now()))

	list, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestVenueRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewVenueRepository(db)
	id := 99

	mock.ExpectQuery("SELECT .* FROM venues WHERE id = \\$1").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	// The repository returns nil, nil or specific error depending on implementation.
	// Looking at implementation: it returns (nil, &apperrors.NotFoundError)
	_, err = repo.GetByID(context.Background(), id)
	assert.Error(t, err)
}

func TestVenueRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewVenueRepository(db)
	venue := &model.Venue{ID: 1, Name: "Venue Updated", Location: "Loc Updated", Description: "Desc Updated"}

	// Success case
	mock.ExpectExec("UPDATE venues SET").
		WithArgs(venue.Name, venue.Location, venue.Description, venue.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), venue)
	assert.NoError(t, err)

	// Not Found case
	mock.ExpectExec("UPDATE venues SET").
		WithArgs(venue.Name, venue.Location, venue.Description, venue.ID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err = repo.Update(context.Background(), venue)
	assert.Error(t, err)
	// Optionally check if error is NotFoundError if you want to be precise
}

func TestVenueRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewVenueRepository(db)
	id := 1

	// Success case
	mock.ExpectExec("DELETE FROM venues WHERE id =").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), id)
	assert.NoError(t, err)

	// Not Found case
	mock.ExpectExec("DELETE FROM venues WHERE id =").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err = repo.Delete(context.Background(), id)
	assert.Error(t, err)
}
