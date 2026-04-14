package job

import (
	"context"
)

// Handler defines the interface for processing an asynchronous job.
type Handler interface {
	Handle(ctx context.Context, payload []byte) error
}
