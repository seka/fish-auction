package entity

import "github.com/seka/fish-auction/backend/internal/domain/model"

type Fisherman struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (e *Fisherman) ToModel() *model.Fisherman {
	return &model.Fisherman{
		ID:   e.ID,
		Name: e.Name,
	}
}
