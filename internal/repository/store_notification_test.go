package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestClaimNotifications(t *testing.T) {
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

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "order_id", "order_status", "attempts", "user_email"}).
					AddRow(1, 10, "shipped", 1, "younes@example.com")

				s.ExpectQuery(query).WithArgs(50, float64(120)).WillReturnRows(rows)

				got, err := ps.ClaimNotifications(context.Background(), 50, 2*time.Minute)
				require.NoError(t, err)

				require.Equal(t, []*model.Notification{{
					ID:          1,
					OrderID:     10,
					OrderStatus: model.OrderShipped,
					UserEmail:   "younes@example.com",
					Attempts:    1,
				}}, got)

				require.NoError(t, s.ExpectationsWereMet())
			},
		},
		{
			name: "failed claiming",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectQuery(query).WithArgs(50, float64(120)).WillReturnError(fmt.Errorf("error claiming"))

				_, err := ps.ClaimNotifications(context.Background(), 50, 2*time.Minute)
				require.Error(t, err)

				require.NoError(t, s.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
				tc.test(t, NewPostgresStorer(db), s)
			})
		})
	}
}

func TestMarkNotificationSent(t *testing.T) {
	query := `
		UPDATE notifications
		SET state = 'sent', sent_at = NOW(), locked_until = NULL, last_error = NULL
		WHERE id = $1
	`

	withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
		s.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		err := NewPostgresStorer(db).MarkNotificationSent(context.Background(), 1)
		require.NoError(t, err)

		require.NoError(t, s.ExpectationsWereMet())
	})
}

func TestMarkNotificationFailed(t *testing.T) {
	query := `
		UPDATE notifications
		SET last_error = $2,
			state = CASE WHEN attempts >= $3 THEN 'failed' ELSE 'pending' END
		WHERE id = $1
	`

	withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
		s.ExpectExec(query).WithArgs(1, "mailbox unavailable", 5).WillReturnResult(sqlmock.NewResult(0, 1))

		err := NewPostgresStorer(db).MarkNotificationFailed(context.Background(), 1, "mailbox unavailable", 5)
		require.NoError(t, err)

		require.NoError(t, s.ExpectationsWereMet())
	})
}
