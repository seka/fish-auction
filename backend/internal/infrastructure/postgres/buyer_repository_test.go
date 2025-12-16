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

type mockBuyerCache struct{}

func (m *mockBuyerCache) Get(ctx context.Context, id int) (*model.Buyer, error) {
	return nil, errors.New("cache miss")
}
func (m *mockBuyerCache) Set(ctx context.Context, id int, buyer *model.Buyer) error { return nil }
func (m *mockBuyerCache) Delete(ctx context.Context, id int) error                  { return nil }

func TestBuyerRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewBuyerRepository(db, &mockBuyerCache{})
	buyer := &model.Buyer{Name: "Buyer1", Organization: "Org1", ContactInfo: "Contact1"}

	mock.ExpectQuery("INSERT INTO buyers").
		WithArgs(buyer.Name, buyer.Organization, buyer.ContactInfo).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	created, err := repo.Create(context.Background(), buyer)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
}

func TestBuyerRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewBuyerRepository(db, &mockBuyerCache{})
	id := 1

	mock.ExpectQuery("SELECT id, name, organization, contact_info FROM buyers WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "organization", "contact_info"}).
			AddRow(1, "Buyer1", "Org1", "Contact1"))

	found, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, found.ID)
}
