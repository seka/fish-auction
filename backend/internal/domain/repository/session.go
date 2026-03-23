package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// SessionRepository provides SessionRepository related functionality.
type SessionRepository interface {
	Create(ctx context.Context, userID int, role model.SessionRole) (string, error)
	FindByID(ctx context.Context, sessionID string) (*model.Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteAllByUserID(ctx context.Context, userID int, role model.SessionRole) error
}
