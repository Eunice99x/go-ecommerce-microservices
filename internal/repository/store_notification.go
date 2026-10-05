package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
)

// ClaimNotifications locks rows with a lease so two workers never send the same email;
// if a worker dies the lease expires and the row is retried
func (ps *PostgresStorer) ClaimNotifications(ctx context.Context, limit int, lease time.Duration) ([]*model.Notification, error) {
	query := `
		UPDATE notifications n
		SET attempts = n.attempts + 1, locked_until = NOW() + make_interval(secs => $2)
		FROM orders o, users u
		WHERE n.id IN (
			SELECT id FROM notifications
			WHERE state = 'pending' AND (locked_until IS NULL OR locked_until < NOW())
			ORDER BY created_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		AND o.id = n.order_id
		AND u.id = o.user_id
		RETURNING n.id, n.order_id, n.order_status, n.attempts, u.email AS user_email
	`

	var ns []*model.Notification

	err := ps.db.SelectContext(ctx, &ns, query, limit, lease.Seconds())
	if err != nil {
		return nil, fmt.Errorf("error claiming notifications: %w", err)
	}

	return ns, nil
}

func (ps *PostgresStorer) MarkNotificationSent(ctx context.Context, id int64) error {
	query := `
		UPDATE notifications
		SET state = 'sent', sent_at = NOW(), locked_until = NULL, last_error = NULL
		WHERE id = $1
	`

	_, err := ps.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error marking notification sent: %w", err)
	}

	return nil
}

// MarkNotificationFailed keeps the lease, so it also works as the retry delay
func (ps *PostgresStorer) MarkNotificationFailed(ctx context.Context, id int64, reason string, maxAttempts int) error {
	query := `
		UPDATE notifications
		SET last_error = $2,
			state = CASE WHEN attempts >= $3 THEN 'failed' ELSE 'pending' END
		WHERE id = $1
	`

	_, err := ps.db.ExecContext(ctx, query, id, reason, maxAttempts)
	if err != nil {
		return fmt.Errorf("error marking notification failed: %w", err)
	}

	return nil
}
