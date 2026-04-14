package worker

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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

// Dispatch routes the message to the correct handler based on the JobType message attribute.
func (d *Dispatcher) Dispatch(ctx context.Context, msg types.Message) error {
	jobTypeStr := ""
	if attr, ok := msg.MessageAttributes["JobType"]; ok {
		jobTypeStr = *attr.StringValue
	}

	jobType := model.JobType(jobTypeStr)
	handler, ok := d.handlers[jobType]
	if !ok {
		return fmt.Errorf("no handler registered for job type: %s", jobType)
	}

	return handler.Handle(ctx, []byte(*msg.Body))
}
