package entity

import (
	"strings"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Fisherman struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (e *Fisherman) Validate() error {
	if strings.TrimSpace(e.Name) == "" {
		return &errors.ValidationError{
			Field:   "name",
			Message: "cannot be empty",
		}
	}
	return nil
}

func (e *Fisherman) ToModel() *model.Fisherman {
	return &model.Fisherman{
		ID:   e.ID,
		Name: e.Name,
	}
}
