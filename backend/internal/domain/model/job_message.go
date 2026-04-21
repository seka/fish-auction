package model

// JobMessage represents a generic message from a job queue.
type JobMessage struct {
	ID            string
	ReceiptHandle string // Implementation-specific handle for deletion
	JobType       JobType
	Payload       []byte
	ReceiveCount  int // Number of times this message has been received
}
