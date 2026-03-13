package postgres

import (
	"context"
	"database/sql"
	"errors"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type buyerStore struct {
	db datastore.Database
}

var _ repository.BuyerRepository = (*buyerStore)(nil)

// NewBuyerStore creates a new instance of BuyerRepository
func NewBuyerStore(db datastore.Database) *buyerStore {
	return &buyerStore{db: db}
}

func (r *buyerStore) Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) { // Changed receiver type
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
		return nil, err
	}
	buyer.ID = e.ID
	return buyer, nil
}

func (r *buyerStore) List(ctx context.Context) ([]model.Buyer, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, organization, contact_info FROM buyers WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
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
	return buyers, rows.Err()
}

func (r *buyerStore) FindByID(ctx context.Context, id int) (*model.Buyer, error) {
	var e entity.Buyer
	err := r.db.QueryRow(ctx,
		"SELECT id, name, organization, contact_info FROM buyers WHERE id = $1",
		id,
	).Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Buyer", ID: id}
		}
		return nil, err
	}

	return e.ToModel(), nil
}

func (r *buyerStore) FindByName(ctx context.Context, name string) (*model.Buyer, error) {
	var e entity.Buyer
	err := r.db.QueryRow(ctx,
		"SELECT id, name, organization, contact_info FROM buyers WHERE name = $1 AND deleted_at IS NULL",
		name,
	).Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Buyer", ID: 0} // ID unknown
		}
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *buyerStore) FindByEmail(ctx context.Context, email string) (*model.Buyer, error) {
	var e entity.Buyer
	query := `
		SELECT b.id, b.name, b.organization, b.contact_info
		FROM buyers b
		JOIN authentications a ON b.id = a.buyer_id
		WHERE a.email = $1 AND b.deleted_at IS NULL
	`
	err := r.db.QueryRow(ctx, query, email).Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No buyer found with this email
			return nil, &apperrors.NotFoundError{Resource: "Buyer", ID: 0}
		}
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *buyerStore) Delete(ctx context.Context, id int) error {
	_, err := r.db.Execute(ctx, "UPDATE buyers SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", id)
	return err
}
