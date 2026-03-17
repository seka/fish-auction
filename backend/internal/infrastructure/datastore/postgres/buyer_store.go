package postgres

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

// BuyerStore implements repository.BuyerRepository using PostgreSQL.
type BuyerStore struct {
	db datastore.Database
}

var _ repository.BuyerRepository = (*BuyerStore)(nil)

// NewBuyerStore creates a new instance of BuyerRepository
func NewBuyerStore(db datastore.Database) *BuyerStore {
	return &BuyerStore{db: db}
}

// Create stores a new buyer.
func (r *BuyerStore) Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
	e := entity.Buyer{
		Name:         buyer.Name,
		Organization: buyer.Organization,
		ContactInfo:  buyer.ContactInfo,
	}
	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRow(ctx,
		"INSERT INTO buyers (name, organization, contact_info) VALUES ($1, $2, $3) RETURNING id",
		e.Name, e.Organization, e.ContactInfo).Scan(&e.ID)
	if err != nil {
		return nil, dserrors.HandleError(err, "Buyer", 0, "Create")
	}
	buyer.ID = e.ID
	return buyer, nil
}

// List returns all active buyers.
func (r *BuyerStore) List(ctx context.Context) ([]model.Buyer, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, organization, contact_info FROM buyers WHERE deleted_at IS NULL")
	if err != nil {
		return nil, dserrors.HandleError(err, "Buyer", 0, "List")
	}
	defer func() { _ = rows.Close() }()

	var buyers []model.Buyer
	for rows.Next() {
		var e entity.Buyer
		if err := rows.Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo); err != nil {
			return nil, err
		}
		buyers = append(buyers, *e.ToModel())
	}
	return buyers, dserrors.HandleError(rows.Err(), "Buyer", 0, "List")
}

// FindByID returns a buyer by its ID.
func (r *BuyerStore) FindByID(ctx context.Context, id int) (*model.Buyer, error) {
	var e entity.Buyer
	err := r.db.QueryRow(ctx,
		"SELECT id, name, organization, contact_info FROM buyers WHERE id = $1",
		id,
	).Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo)
	if err != nil {
		return nil, dserrors.HandleError(err, "Buyer", id, "FindByID")
	}

	return e.ToModel(), nil
}

// FindByName returns a buyer by its name.
func (r *BuyerStore) FindByName(ctx context.Context, name string) (*model.Buyer, error) {
	var e entity.Buyer
	err := r.db.QueryRow(ctx,
		"SELECT id, name, organization, contact_info FROM buyers WHERE name = $1 AND deleted_at IS NULL",
		name,
	).Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo)
	if err != nil {
		return nil, dserrors.HandleError(err, "Buyer", 0, "FindByName")
	}
	return e.ToModel(), nil
}

// FindByEmail returns a buyer by its authentication email.
func (r *BuyerStore) FindByEmail(ctx context.Context, email string) (*model.Buyer, error) {
	var e entity.Buyer
	query := `
		SELECT b.id, b.name, b.organization, b.contact_info
		FROM buyers b
		JOIN authentications a ON b.id = a.buyer_id
		WHERE a.email = $1 AND b.deleted_at IS NULL
	`
	err := r.db.QueryRow(ctx, query, email).Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo)
	if err != nil {
		return nil, dserrors.HandleError(err, "Buyer", 0, "FindByEmail")
	}
	return e.ToModel(), nil
}

// Delete marks a buyer as deleted.
func (r *BuyerStore) Delete(ctx context.Context, id int) error {
	_, err := r.db.Execute(ctx, "UPDATE buyers SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", id)
	if err != nil {
		return dserrors.HandleError(err, "Buyer", id, "Delete")
	}
	return nil
}
