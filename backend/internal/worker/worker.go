package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

const (
	waitTimeSeconds = 20
	shutdownTimeout = 30 * time.Second
	retryDelay      = 5 * time.Second
)

// queuePoller is a helper interface for runLoop to poll any domain queue.
type queuePoller interface {
	Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error)
	DeleteMessage(ctx context.Context, message *model.JobMessage) error
}

// HandlerFunc is a function that processes a job message.
type HandlerFunc func(ctx context.Context, msg *model.JobMessage) error

// Worker represents the background job worker.
type Worker struct {
	queue        queuePoller
	emailHandler HandlerFunc
	pushHandler  HandlerFunc
	wg           sync.WaitGroup
}

// NewWorker creates a new Worker instance.
func NewWorker(
	queue queuePoller,
	emailHandler HandlerFunc,
	pushHandler HandlerFunc,
) *Worker {
	return &Worker{
		queue:        queue,
		emailHandler: emailHandler,
		pushHandler:  pushHandler,
	}
}

// Start runs the worker polling loops and blocks until the context is canceled.
func (w *Worker) Start(ctx context.Context) error {
	log.Println("Worker starting...")

	// Start a single polling loop and dispatch by JobType.
	w.wg.Add(1)
	go w.runLoop(ctx, w.queue, "JobQueue")

	// Block until context is canceled
	<-ctx.Done()
	log.Println("Worker received shutdown signal. Shutting down gracefully...")

	// Graceful shutdown wait
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Worker finished all jobs")
	case <-shutdownCtx.Done():
		log.Println("Worker shutdown timed out, some jobs may have been interrupted")
	}

	log.Println("Worker exiting")
	return nil
}

func (w *Worker) runLoop(ctx context.Context, poller queuePoller, name string) {
	defer w.wg.Done()
	log.Printf("Worker: starting polling loop for %s queue...", name)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			messages, err := poller.Dequeue(ctx, waitTimeSeconds)
			if err != nil {
				// Avoid log spamming if the context is canceled
				if ctx.Err() != nil {
					return
				}
				log.Printf("Worker (%s): error receiving messages: %v", name, err)
				time.Sleep(retryDelay) // Wait before retrying
				continue
			}

			for _, msg := range messages {
				handler, err := w.selectHandler(msg.JobType)
				if err != nil {
					log.Printf("Worker (%s): unsupported job type for message %v: %v", name, msg.ID, err)
					continue
				}

				if err := handler(ctx, msg); err != nil {
					// NOTE: 処理失敗時はメッセージを削除せず、SQS の Visibility Timeout 後の再配信に任せます。
					log.Printf("Worker (%s): error processing message %v (attempt %d): %v", name, msg.ID, msg.ReceiveCount, err)
					continue
				}

				if err := poller.DeleteMessage(ctx, msg); err != nil {
					log.Printf("Worker (%s): error deleting message %v: %v", name, msg.ID, err)
				}
			}
		}
	}
}

func (w *Worker) selectHandler(jobType model.JobType) (HandlerFunc, error) {
	switch jobType {
	case model.JobTypeEmail:
		return w.emailHandler, nil
	case model.JobTypePushNotification:
		return w.pushHandler, nil
	default:
		return nil, fmt.Errorf("unsupported job type: %s", jobType)
	}
}
