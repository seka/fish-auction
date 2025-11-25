package entity

import (
	"strings"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type AuctionItem struct {
	ID          int       `db:"id"`
	FishermanID int       `db:"fisherman_id"`
	FishType    string    `db:"fish_type"`
	Quantity    int       `db:"quantity"`
	Unit        string    `db:"unit"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}

func (e *AuctionItem) Validate() error {
	if e.FishermanID <= 0 {
		return &errors.ValidationError{
			Field:   "fisherman_id",
			Message: "must be positive",
		}
	}
	if strings.TrimSpace(e.FishType) == "" {
		return &errors.ValidationError{
			Field:   "fish_type",
			Message: "cannot be empty",
		}
	}
	if e.Quantity <= 0 {
		return &errors.ValidationError{
			Field:   "quantity",
			Message: "must be positive",
		}
	}
	if strings.TrimSpace(e.Unit) == "" {
		return &errors.ValidationError{
			Field:   "unit",
			Message: "cannot be empty",
		}
	}
	return nil
}

func (e *AuctionItem) ToModel() *model.AuctionItem {
	return &model.AuctionItem{
		ID:          e.ID,
		FishermanID: e.FishermanID,
		FishType:    e.FishType,
		Quantity:    e.Quantity,
		Unit:        e.Unit,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
	}
}
