package entity

import "github.com/seka/fish-auction/backend/internal/domain/model"

type Buyer struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (e *Buyer) ToModel() *model.Buyer {
	return &model.Buyer{
		ID:   e.ID,
		Name: e.Name,
	}
}
