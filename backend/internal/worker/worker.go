package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
)

const (
	shutdownTimeout = 30 * time.Second
	retryDelay      = 5 * time.Second
)

// HandlerFunc is a function that processes a job message.
type HandlerFunc func(ctx context.Context, msg *model.JobMessage) error

// Worker represents the background job worker.
type Worker struct {
	queue        service.JobQueue
	emailHandler HandlerFunc
	pushHandler  HandlerFunc
	waitTime     int32
	wg           sync.WaitGroup
}

// NewWorker creates a new Worker instance.
func NewWorker(
	queue service.JobQueue,
	emailHandler HandlerFunc,
	pushHandler HandlerFunc,
	waitTime int32,
) *Worker {
	return &Worker{
		queue:        queue,
		emailHandler: emailHandler,
		pushHandler:  pushHandler,
		waitTime:     waitTime,
	}
}

// Start runs the worker polling loops and blocks until the context is canceled.
func (w *Worker) Start(ctx context.Context) error {
	log.Println("Worker starting...")

	// Start a single polling loop and dispatch by JobType.
	w.wg.Add(1)
	go w.runLoop(ctx, w.queue)

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

func (w *Worker) runLoop(ctx context.Context, poller service.JobQueue) {
	defer w.wg.Done()
	log.Println("Worker: starting polling loop...")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			messages, err := poller.Dequeue(ctx, w.waitTime)
			if err != nil {
				// Avoid log spamming if the context is canceled
				if ctx.Err() != nil {
					return
				}
				log.Printf("Worker: error receiving messages: %v", err)
				time.Sleep(retryDelay) // Wait before retrying
				continue
			}

			for _, msg := range messages {
				handler, err := w.selectHandler(msg.JobType)
				if err != nil {
					log.Printf("Worker: unsupported job type for message %v: %v", msg.ID, err)
					continue
				}

				if err := handler(ctx, msg); err != nil {
					// NOTE: 処理失敗時はメッセージを削除せず、SQS の Visibility Timeout 後の再配信に任せます。
					// 無限ループを防ぐため、インフラ（SQS）側で DLQ（Dead Letter Queue）および
					// RedrivePolicy（maxReceiveCount）が設定されている必要があります。
					log.Printf("Worker: error processing message %v (attempt %d): %v", msg.ID, msg.ReceiveCount, err)
					continue
				}

				if err := poller.DeleteMessage(ctx, msg); err != nil {
					log.Printf("Worker: error deleting message %v: %v", msg.ID, err)
				}
			}
		}
	}
}

func (w *Worker) selectHandler(jobType model.JobType) (HandlerFunc, error) {
	switch jobType {
	case model.JobTypeEmail:
		return w.emailHandler, nil
	case model.JobTypePushOutbid, model.JobTypePushAuctionStatusChanged:
		return w.pushHandler, nil
	default:
		return nil, fmt.Errorf("unsupported job type: %s", jobType)
	}
}
