package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

const sessionKeyPrefix = "session:"

type SessionStore struct {
	cache datastore.Cache
	ttl   time.Duration
}

var _ repository.SessionRepository = (*SessionStore)(nil)

func NewSessionStore(cache datastore.Cache, ttl time.Duration) *SessionStore {
	return &SessionStore{
		cache: cache,
		ttl:   ttl,
	}
}

func (s *SessionStore) Create(ctx context.Context, userID int, role model.SessionRole) (string, error) {
	sessionID, err := newSessionID()
	if err != nil {
		return "", err
	}

	session := model.Session{
		ID:        sessionID,
		UserID:    userID,
		Role:      role,
		CreatedAt: time.Now().UTC(),
	}

	payload, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("marshal session: %w", err)
	}

	if err := s.cache.Set(ctx, sessionKey(sessionID), payload, s.ttl); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *SessionStore) FindByID(ctx context.Context, sessionID string) (*model.Session, error) {
	payload, err := s.cache.Get(ctx, sessionKey(sessionID))
	if err != nil {
		return nil, err
	}
	if payload == nil {
		return nil, nil
	}

	var session model.Session
	if err := json.Unmarshal(payload, &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &session, nil
}

func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}

	return s.cache.Delete(ctx, sessionKey(sessionID))
}

func sessionKey(sessionID string) string {
	return sessionKeyPrefix + sessionID
}

func newSessionID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}

	return hex.EncodeToString(buf), nil
}
