package testing

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type MockSessionRepository struct {
	CreateFunc         func(ctx context.Context, userID int, role model.SessionRole) (string, error)
	FindByIDFunc       func(ctx context.Context, sessionID string) (*model.Session, error)
	DeleteFunc         func(ctx context.Context, sessionID string) error
	Sessions           map[string]*model.Session
	NextSessionID      string
	DeletedSessionIDs  []string
}

func (m *MockSessionRepository) Create(ctx context.Context, userID int, role model.SessionRole) (string, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, userID, role)
	}

	sessionID := m.NextSessionID
	if sessionID == "" {
		sessionID = "test-session"
	}
	if m.Sessions == nil {
		m.Sessions = make(map[string]*model.Session)
	}

	m.Sessions[sessionID] = &model.Session{
		ID:     sessionID,
		UserID: userID,
		Role:   role,
	}

	return sessionID, nil
}

func (m *MockSessionRepository) FindByID(ctx context.Context, sessionID string) (*model.Session, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, sessionID)
	}
	if m.Sessions == nil {
		return nil, nil
	}
	return m.Sessions[sessionID], nil
}

func (m *MockSessionRepository) Delete(ctx context.Context, sessionID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, sessionID)
	}

	m.DeletedSessionIDs = append(m.DeletedSessionIDs, sessionID)
	if m.Sessions != nil {
		delete(m.Sessions, sessionID)
	}

	return nil
}

func (m *MockSessionRepository) DeleteAllByUserID(ctx context.Context, userID int, role model.SessionRole) error {
	if m.Sessions == nil {
		return nil
	}

	for id, s := range m.Sessions {
		if s.UserID == userID && s.Role == role {
			delete(m.Sessions, id)
			m.DeletedSessionIDs = append(m.DeletedSessionIDs, id)
		}
	}

	return nil
}
