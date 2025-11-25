package entity

import (
	"errors"
	"strings"
	"time"

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
		return errors.New("fisherman_id must be positive")
	}
	if strings.TrimSpace(e.FishType) == "" {
		return errors.New("fish_type cannot be empty")
	}
	if e.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if strings.TrimSpace(e.Unit) == "" {
		return errors.New("unit cannot be empty")
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
