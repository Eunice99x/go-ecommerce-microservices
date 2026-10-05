package notifier

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/stretchr/testify/require"
)

type fakeQueue struct {
	notifications []*model.Notification
	claimErr      error

	mu        sync.Mutex
	completed map[int64]error
}

func (f *fakeQueue) ClaimNotifications(ctx context.Context, limit int) ([]*model.Notification, error) {
	return f.notifications, f.claimErr
}

func (f *fakeQueue) CompleteNotification(ctx context.Context, id int64, sendErr error) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.completed == nil {
		f.completed = map[int64]error{}
	}
	f.completed[id] = sendErr

	return nil
}

// fakeMailer fails for the addresses in failFor and records the most sends running at once
type fakeMailer struct {
	failFor map[string]bool
	delay   time.Duration

	mu          sync.Mutex
	inFlight    int
	maxInFlight int
}

func (f *fakeMailer) Send(ctx context.Context, to, subject, body string) error {
	f.mu.Lock()
	f.inFlight++
	f.maxInFlight = max(f.maxInFlight, f.inFlight)
	f.mu.Unlock()

	defer func() {
		f.mu.Lock()
		f.inFlight--
		f.mu.Unlock()
	}()

	time.Sleep(f.delay)

	if f.failFor[to] {
		return errors.New("mailbox unavailable")
	}

	return nil
}

func notifications(n int) []*model.Notification {
	ns := make([]*model.Notification, 0, n)
	for i := 1; i <= n; i++ {
		ns = append(ns, &model.Notification{
			ID:          int64(i),
			OrderID:     int64(i),
			OrderStatus: model.OrderShipped,
			UserEmail:   fmt.Sprintf("user%d@example.com", i),
		})
	}

	return ns
}

func TestProcessBatch(t *testing.T) {
	t.Run("completes every notification and reports failures", func(t *testing.T) {
		q := &fakeQueue{notifications: notifications(3)}
		m := &fakeMailer{failFor: map[string]bool{"user2@example.com": true}}

		w := NewWorker(q, m, Config{BatchSize: 10, MaxConcurrent: 2})

		require.NoError(t, w.processBatch(t.Context()))

		require.Len(t, q.completed, 3)
		require.NoError(t, q.completed[1])
		require.Error(t, q.completed[2])
		require.NoError(t, q.completed[3])
	})

	t.Run("never exceeds max concurrent sends", func(t *testing.T) {
		q := &fakeQueue{notifications: notifications(20)}
		m := &fakeMailer{delay: 10 * time.Millisecond}

		w := NewWorker(q, m, Config{BatchSize: 20, MaxConcurrent: 3})

		require.NoError(t, w.processBatch(t.Context()))

		require.Len(t, q.completed, 20)
		require.LessOrEqual(t, m.maxInFlight, 3)
		require.Greater(t, m.maxInFlight, 1, "sends should actually run concurrently")
	})

	t.Run("cancelled context still finishes started sends", func(t *testing.T) {
		q := &fakeQueue{notifications: notifications(5)}
		m := &fakeMailer{delay: 20 * time.Millisecond}

		w := NewWorker(q, m, Config{BatchSize: 5, MaxConcurrent: 2})

		ctx, cancel := context.WithCancel(t.Context())
		time.AfterFunc(5*time.Millisecond, cancel)

		require.NoError(t, w.processBatch(ctx))

		// the first two were started before the cancel and must be completed, not abandoned;
		// whatever wasn't started stays claimed and is retried after the lease expires
		require.GreaterOrEqual(t, len(q.completed), 2)
		for id, err := range q.completed {
			require.NoError(t, err, "notification %d", id)
		}
	})

	t.Run("claim error", func(t *testing.T) {
		q := &fakeQueue{claimErr: errors.New("grpc unavailable")}

		w := NewWorker(q, &fakeMailer{}, Config{BatchSize: 5, MaxConcurrent: 2})

		require.Error(t, w.processBatch(t.Context()))
		require.Empty(t, q.completed)
	})
}
