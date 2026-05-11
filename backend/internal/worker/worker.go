package worker

import (
	"context"
	"fmt"
	"log/slog"
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
	logger       *slog.Logger
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
		logger:       slog.With("component", "worker"),
	}
}

// Start runs the worker polling loops and blocks until the context is canceled.
func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info("worker starting")

	// Start a single polling loop and dispatch by JobType.
	w.wg.Add(1)
	go w.runLoop(ctx, w.queue)

	// Block until context is canceled
	<-ctx.Done()
	w.logger.Info("worker received shutdown signal; shutting down gracefully")

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
		w.logger.Info("worker finished all jobs")
	case <-shutdownCtx.Done():
		w.logger.Warn("worker shutdown timed out; some jobs may have been interrupted")
	}

	w.logger.Info("worker exiting")
	return nil
}

func (w *Worker) runLoop(ctx context.Context, poller service.JobQueue) {
	defer w.wg.Done()
	w.logger.Info("starting polling loop")

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
				w.logger.Error("error receiving messages", "err", err)
				time.Sleep(retryDelay) // Wait before retrying
				continue
			}

			for _, msg := range messages {
				handler, err := w.selectHandler(msg.JobType)
				if err != nil {
					w.logger.Warn("unsupported job type", "message_id", msg.ID, "err", err)
					continue
				}

				if err := handler(ctx, msg); err != nil {
					// NOTE: 処理失敗時はメッセージを削除せず、SQS の Visibility Timeout 後の再配信に任せます。
					// 無限ループを防ぐため、インフラ（SQS）側で DLQ（Dead Letter Queue）および
					// RedrivePolicy（maxReceiveCount）が設定されている必要があります。
					w.logger.Error("error processing message", "message_id", msg.ID, "attempt", msg.ReceiveCount, "err", err)
					continue
				}

				if err := poller.DeleteMessage(ctx, msg); err != nil {
					w.logger.Error("error deleting message", "message_id", msg.ID, "err", err)
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
