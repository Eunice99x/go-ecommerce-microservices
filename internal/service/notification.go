package service

import (
	"context"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
)

const (
	// how long a worker owns a claimed notification, also the retry delay
	notificationLease = 2 * time.Minute

	maxNotificationAttempts = 5
)

func (s *Service) ClaimNotifications(ctx context.Context, limit int) ([]*model.Notification, error) {
	return s.storer.ClaimNotifications(ctx, limit, notificationLease)
}

func (s *Service) MarkNotificationSent(ctx context.Context, id int64) error {
	return s.storer.MarkNotificationSent(ctx, id)
}

func (s *Service) MarkNotificationFailed(ctx context.Context, id int64, reason string) error {
	return s.storer.MarkNotificationFailed(ctx, id, reason, maxNotificationAttempts)
}
