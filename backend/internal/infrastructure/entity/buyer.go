package entity

import (
	"strings"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Buyer struct {
	ID           int        `db:"id"`
	Name         string     `db:"name"`
	Organization string     `db:"organization"`
	ContactInfo  string     `db:"contact_info"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

func (b *Buyer) Validate() error {
	if strings.TrimSpace(b.Name) == "" {
		return &apperrors.ValidationError{Field: "name", Message: "Name is required"}
	}
	if strings.TrimSpace(b.Organization) == "" {
		return &apperrors.ValidationError{Field: "organization", Message: "Organization is required"}
	}
	if strings.TrimSpace(b.ContactInfo) == "" {
		return &apperrors.ValidationError{Field: "contact_info", Message: "Contact info is required"}
	}
	return nil
}

func (b *Buyer) ToModel() *model.Buyer {
	return &model.Buyer{
		ID:           b.ID,
		Name:         b.Name,
		Organization: b.Organization,
		ContactInfo:  b.ContactInfo,
	}
}
