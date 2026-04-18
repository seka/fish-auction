package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/service"
)

const (
	waitTimeSeconds = 20
	shutdownTimeout = 30 * time.Second
	retryDelay      = 5 * time.Second
)

// Worker represents the background job worker.
type Worker struct {
	queue      service.JobQueue
	dispatcher *Dispatcher
	wg         sync.WaitGroup
}

// NewWorker creates a new Worker instance.
func NewWorker(queue service.JobQueue, dispatcher *Dispatcher) *Worker {
	return &Worker{
		queue:      queue,
		dispatcher: dispatcher,
	}
}

// Start runs the worker polling loop and blocks until the context is canceled.
func (w *Worker) Start(ctx context.Context) error {
	log.Println("Worker starting...")

	// Run polling loop in a separate goroutine
	w.wg.Add(1)
	go w.runLoop(ctx)

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

func (w *Worker) runLoop(ctx context.Context) {
	defer w.wg.Done()
	log.Println("Worker starting polling loop...")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			messages, err := w.queue.Dequeue(ctx, waitTimeSeconds)
			if err != nil {
				// Avoid log spamming if the context is cancelled
				if ctx.Err() != nil {
					return
				}
				log.Printf("Error receiving messages: %v", err)
				time.Sleep(retryDelay) // Wait before retrying
				continue
			}

			for _, msg := range messages {
				if err := w.dispatcher.Dispatch(ctx, msg); err != nil {
					// Message is not deleted, will be retried after visibility timeout
					log.Printf("Error processing message %v: %v", msg.ID, err)
					continue
				}

				if err := w.queue.DeleteMessage(ctx, msg); err != nil {
					log.Printf("Error deleting message %v: %v", msg.ID, err)
				}
			}
		}
	}
}
