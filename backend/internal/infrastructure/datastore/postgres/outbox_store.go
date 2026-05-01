package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/event"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

// OutboxStore implements repository.OutboxRepository.
type OutboxStore struct {
	db datastore.Database
}

var _ repository.OutboxRepository = (*OutboxStore)(nil)

// NewOutboxStore creates a new OutboxStore instance.
func NewOutboxStore(db datastore.Database) *OutboxStore {
	return &OutboxStore{db: db}
}

func (s *OutboxStore) insert(ctx context.Context, jobType model.JobType, schemaVersion int, payload []byte) error {
	query := `
		INSERT INTO outbox (job_type, schema_version, payload)
		VALUES ($1, $2, $3)
	`
	if _, err := s.db.Execute(ctx, query, string(jobType), schemaVersion, payload); err != nil {
		return fmt.Errorf("failed to insert outbox message: %w", err)
	}
	return nil
}

// InsertEmailJob serializes and inserts an email job.
func (s *OutboxStore) InsertEmailJob(ctx context.Context, to, resetURL, emailType string) error {
	msg := event.EmailMessage{
		EmailType: event.EmailType(emailType),
		To:        to,
		ResetURL:  resetURL,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal email job: %w", err)
	}
	return s.insert(ctx, model.JobTypeEmail, 1, payload)
}

// InsertPushNotificationJob serializes and inserts a push notification job.
func (s *OutboxStore) InsertPushNotificationJob(ctx context.Context, buyerID int, payload any) error {
	msg := event.PushNotificationMessage{
		BuyerID: buyerID,
		Payload: payload,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal push notification job: %w", err)
	}
	return s.insert(ctx, model.JobTypePushNotification, 1, body)
}

func (s *OutboxStore) Claim(ctx context.Context, limit int, claimedBy string) ([]*model.OutboxMessage, error) {
	query := `
		UPDATE outbox
		SET status = 'processing', claimed_at = NOW(), claimed_by = $2, attempts = attempts + 1
		WHERE id IN (
			SELECT id FROM outbox
			WHERE status = 'pending' AND available_at <= NOW()
			ORDER BY id
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, job_type, schema_version, payload, status, attempts, max_attempts, available_at, created_at
	`

	rows, err := s.db.Query(ctx, query, limit, claimedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to claim outbox messages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var msgs []*model.OutboxMessage
	for rows.Next() {
		msg := &model.OutboxMessage{}
		var statusStr string
		var jobTypeStr string
		if err := rows.Scan(
			&msg.ID, &jobTypeStr, &msg.SchemaVersion, &msg.Payload,
			&statusStr, &msg.Attempts, &msg.MaxAttempts, &msg.AvailableAt, &msg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan outbox message: %w", err)
		}
		msg.JobType = model.JobType(jobTypeStr)
		msg.Status = model.OutboxStatus(statusStr)
		msgs = append(msgs, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate outbox messages: %w", err)
	}
	return msgs, nil
}

func (s *OutboxStore) MarkProcessed(ctx context.Context, ids []int64) error {
	query := `UPDATE outbox SET status = 'processed', processed_at = NOW(), claimed_at = NULL, claimed_by = NULL WHERE id = ANY($1)`
	if _, err := s.db.Execute(ctx, query, pq.Array(ids)); err != nil {
		return fmt.Errorf("failed to mark outbox messages as processed: %w", err)
	}
	return nil
}

func (s *OutboxStore) MarkFailed(ctx context.Context, id int64, lastError string) error {
	query := `
		UPDATE outbox
		SET status = CASE WHEN attempts >= max_attempts THEN 'failed' ELSE 'pending' END,
		    available_at = CASE WHEN attempts >= max_attempts THEN available_at
		                        ELSE NOW() + (INTERVAL '1 second' * POWER(2, LEAST(attempts, 10)))
		                   END,
		    last_error = $2,
		    claimed_at = NULL,
		    claimed_by = NULL
		WHERE id = $1
	`
	if _, err := s.db.Execute(ctx, query, id, lastError); err != nil {
		return fmt.Errorf("failed to mark outbox message %d as failed: %w", id, err)
	}
	return nil
}

func (s *OutboxStore) RecoverStale(ctx context.Context, timeout time.Duration) (int64, error) {
	query := `
		UPDATE outbox
		SET status = 'pending', claimed_at = NULL, claimed_by = NULL
		WHERE status = 'processing' AND claimed_at < NOW() - $1::interval
	`
	n, err := s.db.Execute(ctx, query, fmt.Sprintf("%d seconds", int(timeout.Seconds())))
	if err != nil {
		return 0, fmt.Errorf("failed to recover stale outbox messages: %w", err)
	}
	return n, nil
}

func (s *OutboxStore) DeleteProcessedBefore(ctx context.Context, before time.Time) (int64, error) {
	query := `DELETE FROM outbox WHERE status = 'processed' AND processed_at < $1`
	n, err := s.db.Execute(ctx, query, before)
	if err != nil {
		return 0, fmt.Errorf("failed to delete processed outbox messages: %w", err)
	}
	return n, nil
}
