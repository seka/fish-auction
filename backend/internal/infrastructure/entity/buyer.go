package entity

import (
	"strings"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Buyer struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (e *Buyer) Validate() error {
	if strings.TrimSpace(e.Name) == "" {
		return &errors.ValidationError{
			Field:   "name",
			Message: "cannot be empty",
		}
	}
	return nil
}

func (e *Buyer) ToModel() *model.Buyer {
	return &model.Buyer{
		ID:   e.ID,
		Name: e.Name,
	}
}
