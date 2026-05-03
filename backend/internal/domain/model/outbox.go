package model

import "time"

// OutboxStatus represents the processing state of an outbox message.
type OutboxStatus string

const (
	OutboxStatusPending    OutboxStatus = "pending"
	OutboxStatusProcessing OutboxStatus = "processing"
	OutboxStatusProcessed  OutboxStatus = "processed"
	OutboxStatusFailed     OutboxStatus = "failed"
)

// OutboxMessage represents a message stored in the transactional outbox table.
type OutboxMessage struct {
	ID            int64
	JobType       JobType
	SchemaVersion int
	Payload       []byte
	Status        OutboxStatus
	Attempts      int
	MaxAttempts   int
	AvailableAt   time.Time
	CreatedAt     time.Time
}
