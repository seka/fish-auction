package worker

import (
	"context"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/worker/job"
)

// Dispatcher handles job dispatching to the appropriate handler.
type Dispatcher struct {
	handlers map[model.JobType]job.Handler
}

// NewDispatcher creates a new Dispatcher instance.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[model.JobType]job.Handler),
	}
}

// Register adds a handler for a given JobType.
func (d *Dispatcher) Register(jobType model.JobType, handler job.Handler) {
	d.handlers[jobType] = handler
}

// Dispatch routes the message to the correct handler based on the JobType.
func (d *Dispatcher) Dispatch(ctx context.Context, msg *model.JobMessage) error {
	handler, ok := d.handlers[msg.JobType]
	if !ok {
		return fmt.Errorf("no handler registered for job type: %s", msg.JobType)
	}

	return handler.Handle(ctx, msg.Payload)
}
