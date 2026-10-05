// Package notifier sends the queued order emails.
package notifier

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
	"golang.org/x/sync/semaphore"
)

const sendTimeout = 30 * time.Second

// Queue is implemented by the gRPC client
type Queue interface {
	ClaimNotifications(ctx context.Context, limit int) ([]*model.Notification, error)
	CompleteNotification(ctx context.Context, id int64, sendErr error) error
}

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type Config struct {
	Interval      time.Duration // how often to poll for new notifications
	BatchSize     int           // how many to claim per poll
	MaxConcurrent int64         // how many emails can be in flight at once
}

type Worker struct {
	queue  Queue
	mailer Mailer
	cfg    Config
	sem    *semaphore.Weighted
}

func NewWorker(queue Queue, mailer Mailer, cfg Config) *Worker {
	return &Worker{
		queue:  queue,
		mailer: mailer,
		cfg:    cfg,
		sem:    semaphore.NewWeighted(cfg.MaxConcurrent),
	}
}

// Run polls until ctx is cancelled, finishing the current batch first
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		if err := w.processBatch(ctx); err != nil {
			log.Printf("notifier: %v", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// processBatch sends one batch with at most MaxConcurrent emails in flight
func (w *Worker) processBatch(ctx context.Context) error {
	ns, err := w.queue.ClaimNotifications(ctx, w.cfg.BatchSize)
	if err != nil {
		return fmt.Errorf("error claiming notifications: %w", err)
	}

	// started sends should finish even on shutdown
	sendCtx := context.WithoutCancel(ctx)

	var wg sync.WaitGroup

	for _, n := range ns {
		// waits for a free slot; fails only on shutdown, unsent rows get retried after their lease
		if err := w.sem.Acquire(ctx, 1); err != nil {
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer w.sem.Release(1)

			w.deliver(sendCtx, n)
		}()
	}

	wg.Wait()

	return nil
}

func (w *Worker) deliver(ctx context.Context, n *model.Notification) {
	ctx, cancel := context.WithTimeout(ctx, sendTimeout)
	defer cancel()

	subject := fmt.Sprintf("Your order #%d is %s", n.OrderID, n.OrderStatus)
	body := fmt.Sprintf("Hi,\n\nYour order #%d is now %s.\n", n.OrderID, n.OrderStatus)

	sendErr := w.mailer.Send(ctx, n.UserEmail, subject, body)
	if sendErr != nil {
		log.Printf("notifier: sending notification %d (attempt %d): %v", n.ID, n.Attempts, sendErr)
	}

	if err := w.queue.CompleteNotification(ctx, n.ID, sendErr); err != nil {
		log.Printf("notifier: completing notification %d: %v", n.ID, err)
	}
}
