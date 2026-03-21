package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

const sessionKeyPrefix = "session:"

type SessionStore struct {
	cache datastore.Cache
	ttl   time.Duration
}

type sessionJSON struct {
	ID        string            `json:"id"`
	UserID    int               `json:"user_id"`
	Role      model.SessionRole `json:"role"`
	CreatedAt time.Time         `json:"created_at"`
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

	sJSON := sessionJSON{
		ID:        sessionID,
		UserID:    userID,
		Role:      role,
		CreatedAt: time.Now().UTC(),
	}

	payload, err := json.Marshal(sJSON)
	if err != nil {
		return "", fmt.Errorf("marshal session: %w", err)
	}

	if err := s.cache.Set(ctx, sessionKey(sessionID), payload, s.ttl); err != nil {
		return "", err
	}

	// Add to user sessions set
	if rc := s.getRedisClient(); rc != nil {
		setKey := userSessionsKey(role, userID)
		if err := rc.SAdd(ctx, setKey, sessionID).Err(); err != nil {
			return "", fmt.Errorf("add to session set: %w", err)
		}
		if err := rc.Expire(ctx, setKey, s.ttl).Err(); err != nil {
			return "", fmt.Errorf("expire session set: %w", err)
		}
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

	var sJSON sessionJSON
	if err := json.Unmarshal(payload, &sJSON); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &model.Session{
		ID:        sJSON.ID,
		UserID:    sJSON.UserID,
		Role:      sJSON.Role,
		CreatedAt: sJSON.CreatedAt,
	}, nil
}

func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}

	// Find first to get userID and role for set removal
	session, err := s.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if err := s.cache.Delete(ctx, sessionKey(sessionID)); err != nil {
		return err
	}

	if session != nil {
		if rc := s.getRedisClient(); rc != nil {
			_ = rc.SRem(ctx, userSessionsKey(session.Role, session.UserID), sessionID).Err()
		}
	}

	return nil
}

func (s *SessionStore) DeleteAllByUserID(ctx context.Context, userID int, role model.SessionRole) error {
	rc := s.getRedisClient()
	if rc == nil {
		return fmt.Errorf("redis client not available")
	}

	setKey := userSessionsKey(role, userID)
	sessionIDs, err := rc.SMembers(ctx, setKey).Result()
	if err != nil {
		return fmt.Errorf("get user sessions: %w", err)
	}

	for _, id := range sessionIDs {
		_ = s.cache.Delete(ctx, sessionKey(id))
	}

	return rc.Del(ctx, setKey).Err()
}

func (s *SessionStore) getRedisClient() *goredis.Client {
	if c, ok := s.cache.(*Client); ok {
		return c.client
	}
	return nil
}

func userSessionsKey(role model.SessionRole, userID int) string {
	return fmt.Sprintf("user_sessions:%s:%d", role, userID)
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
