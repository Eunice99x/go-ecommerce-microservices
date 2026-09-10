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

func TestCreateSession(t *testing.T) {
	now := time.Now()

	s := &model.Session{
		ID:           "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70",
		UserEmail:    "younes@example.com",
		RefreshToken: "refresh-token",
		IsRevoked:    false,
		ExpiresAt:    now.Add(24 * time.Hour),
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, m sqlmock.Sqlmock) {
				rows := m.NewRows([]string{"created_at"}).AddRow(now)

				m.ExpectQuery(`INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING created_at`).
					WithArgs(
						s.ID,
						s.UserEmail,
						s.RefreshToken,
						s.IsRevoked,
						s.ExpiresAt,
					).
					WillReturnRows(rows)

				got, err := ps.CreateSession(context.Background(), s)

				require.NoError(t, err)
				require.Equal(t, s.ID, got.ID)
				require.Equal(t, now, got.CreatedAt)

				require.NoError(t, m.ExpectationsWereMet())
			},
		},
		{
			name: "failed creating session",
			test: func(t *testing.T, ps *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectQuery(`INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING created_at`).
					WithArgs(
						s.ID,
						s.UserEmail,
						s.RefreshToken,
						s.IsRevoked,
						s.ExpiresAt,
					).
					WillReturnError(fmt.Errorf("db error"))

				_, err := ps.CreateSession(context.Background(), s)

				require.Error(t, err)
				require.NoError(t, m.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, m)
			})
		})
	}
}

func TestGetSession(t *testing.T) {
	now := time.Now()
	id := "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70"

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, m sqlmock.Sqlmock) {
				rows := m.NewRows([]string{
					"id",
					"user_email",
					"refresh_token",
					"is_revoked",
					"created_at",
					"expires_at",
				}).AddRow(
					id,
					"younes@example.com",
					"refresh-token",
					false,
					now,
					now.Add(24*time.Hour),
				)

				m.ExpectQuery("SELECT * FROM sessions WHERE id=$1").
					WithArgs(id).
					WillReturnRows(rows)

				got, err := ps.GetSession(context.Background(), id)

				require.NoError(t, err)
				require.Equal(t, id, got.ID)
				require.Equal(t, "younes@example.com", got.UserEmail)
				require.Equal(t, "refresh-token", got.RefreshToken)
				require.False(t, got.IsRevoked)

				require.NoError(t, m.ExpectationsWereMet())
			},
		},
		{
			name: "failed getting session",
			test: func(t *testing.T, ps *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT * FROM sessions WHERE id=$1").
					WithArgs(id).
					WillReturnError(fmt.Errorf("db error"))

				_, err := ps.GetSession(context.Background(), id)

				require.Error(t, err)
				require.NoError(t, m.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, m)
			})
		})
	}
}

func TestRevokeSession(t *testing.T) {
	id := "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70"

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE sessions SET is_revoked=true WHERE id=$1").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				err := ps.RevokeSession(context.Background(), id)

				require.NoError(t, err)
				require.NoError(t, m.ExpectationsWereMet())
			},
		},
		{
			name: "failed revoking session",
			test: func(t *testing.T, ps *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE sessions SET is_revoked=true WHERE id=$1").
					WithArgs(id).
					WillReturnError(fmt.Errorf("db error"))

				err := ps.RevokeSession(context.Background(), id)

				require.Error(t, err)
				require.NoError(t, m.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, m)
			})
		})
	}
}

func TestDeleteSession(t *testing.T) {
	id := "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70"

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM sessions WHERE id=$1").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				err := ps.DeleteSession(context.Background(), id)

				require.NoError(t, err)
				require.NoError(t, m.ExpectationsWereMet())
			},
		},
		{
			name: "failed deleting session",
			test: func(t *testing.T, ps *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM sessions WHERE id=$1").
					WithArgs(id).
					WillReturnError(fmt.Errorf("db error"))

				err := ps.DeleteSession(context.Background(), id)

				require.Error(t, err)
				require.NoError(t, m.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, m)
			})
		})
	}
}
