package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
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
	adminEmailQueue       service.AdminEmailQueue
	buyerEmailQueue       service.BuyerEmailQueue
	pushNotificationQueue service.PushNotificationQueue
	emailHandler          HandlerFunc
	pushHandler           HandlerFunc
	wg                    sync.WaitGroup
}

// NewWorker creates a new Worker instance.
func NewWorker(
	adminEmailQueue service.AdminEmailQueue,
	buyerEmailQueue service.BuyerEmailQueue,
	pushNotificationQueue service.PushNotificationQueue,
	emailHandler HandlerFunc,
	pushHandler HandlerFunc,
) *Worker {
	return &Worker{
		adminEmailQueue:       adminEmailQueue,
		buyerEmailQueue:       buyerEmailQueue,
		pushNotificationQueue: pushNotificationQueue,
		emailHandler:          emailHandler,
		pushHandler:           pushHandler,
	}
}

// Start runs the worker polling loops and blocks until the context is canceled.
func (w *Worker) Start(ctx context.Context) error {
	log.Println("Worker starting...")

	// Start polling loops for each specialized queue
	w.wg.Add(3)
	go w.runLoop(ctx, w.adminEmailQueue, w.emailHandler, "AdminEmail")
	go w.runLoop(ctx, w.buyerEmailQueue, w.emailHandler, "BuyerEmail")
	go w.runLoop(ctx, w.pushNotificationQueue, w.pushHandler, "PushNotification")

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

func (w *Worker) runLoop(ctx context.Context, poller queuePoller, handler HandlerFunc, name string) {
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
