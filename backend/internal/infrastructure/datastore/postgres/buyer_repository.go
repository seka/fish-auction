package postgres

import (
	"context"
	"database/sql"
	"errors"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	cache "github.com/seka/fish-auction/backend/internal/infrastructure/cache/redis"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type buyerRepository struct {
	db    *sql.DB
	cache cache.BuyerCache
}

func NewBuyerRepository(db *sql.DB, buyerCache cache.BuyerCache) repository.BuyerRepository {
	return &buyerRepository{
		db:    db,
		cache: buyerCache,
	}
}

func (r *buyerRepository) Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error) {
	e := entity.Buyer{
		Name:         buyer.Name,
		Organization: buyer.Organization,
		ContactInfo:  buyer.ContactInfo,
	}
	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO buyers (name, organization, contact_info) VALUES ($1, $2, $3) RETURNING id",
		e.Name, e.Organization, e.ContactInfo).Scan(&e.ID)
	if err != nil {
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *buyerRepository) List(ctx context.Context) ([]model.Buyer, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, organization, contact_info FROM buyers WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buyers []model.Buyer
	for rows.Next() {
		var e entity.Buyer
		if err := rows.Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo); err != nil {
			return nil, err
		}
		buyers = append(buyers, *e.ToModel())
	}
	return buyers, nil
}

func (r *buyerRepository) FindByID(ctx context.Context, id int) (*model.Buyer, error) {
	// キャッシュを確認
	if buyer, err := r.cache.Get(ctx, id); err == nil && buyer != nil {
		return buyer, nil
	}

	// DBから取得
	var e entity.Buyer
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, organization, contact_info FROM buyers WHERE id = $1",
		id,
	).Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Buyer", ID: id}
		}
		return nil, err
	}

	buyer := e.ToModel()

	// キャッシュに保存（エラーは無視）
	_ = r.cache.Set(ctx, id, buyer)

	return buyer, nil
}

func (r *buyerRepository) FindByName(ctx context.Context, name string) (*model.Buyer, error) {
	var e entity.Buyer
	err := r.db.QueryRowContext(ctx,
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

func (r *buyerRepository) FindByEmail(ctx context.Context, email string) (*model.Buyer, error) {
	var e entity.Buyer
	query := `
		SELECT b.id, b.name, b.organization, b.contact_info 
		FROM buyers b 
		JOIN authentications a ON b.id = a.buyer_id 
		WHERE a.email = $1 AND b.deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&e.ID, &e.Name, &e.Organization, &e.ContactInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No buyer found with this email
			return nil, &apperrors.NotFoundError{Resource: "Buyer", ID: 0}
		}
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *buyerRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE buyers SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", id)
	if err != nil {
		return err
	}
	// キャッシュを削除
	_ = r.cache.Delete(ctx, id)
	return nil
}
