package entity

import (
	"errors"
	"strings"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Fisherman struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (e *Fisherman) Validate() error {
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("fisherman name cannot be empty")
	}
	return nil
}

func (e *Fisherman) ToModel() *model.Fisherman {
	return &model.Fisherman{
		ID:   e.ID,
		Name: e.Name,
	}
}
