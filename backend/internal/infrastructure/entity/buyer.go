package entity

import (
	"errors"
	"strings"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Buyer struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (e *Buyer) Validate() error {
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("buyer name cannot be empty")
	}
	return nil
}

func (e *Buyer) ToModel() *model.Buyer {
	return &model.Buyer{
		ID:   e.ID,
		Name: e.Name,
	}
}
