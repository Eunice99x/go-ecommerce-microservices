package repository

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	now := time.Now()

	u := &model.User{
		Name:     "Younes",
		Email:    "younes@example.com",
		Password: "hashed-password",
		IsAdmin:  false,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				rows := s.NewRows([]string{"id", "created_at"}).
					AddRow(1, now)

				s.ExpectQuery(`INSERT INTO users (name, email, password, is_admin) VALUES ($1, $2, $3, $4) RETURNING id, created_at`).
					WithArgs(u.Name, u.Email, u.Password, u.IsAdmin,).WillReturnRows(rows)

				got, err := ps.CreateUser(context.Background(), u)

				require.NoError(t, err)
				require.Equal(t, int64(1), got.ID)
				require.Equal(t, now, got.CreatedAt)

				require.NoError(t, s.ExpectationsWereMet())
			},
		},
		{
			name: "failed creating user",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectQuery(`
					INSERT INTO users (name, email, password, is_admin)
					VALUES ($1, $2, $3, $4)
					RETURNING id, created_at
				`).
					WithArgs(
						u.Name,
						u.Email,
						u.Password,
						u.IsAdmin,
					).
					WillReturnError(fmt.Errorf("db error"))

				_, err := ps.CreateUser(context.Background(), u)

				require.Error(t, err)
				require.NoError(t, s.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, s)
			})
		})
	}
}

func TestGetUser(t *testing.T) {
	now := time.Now()

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				rows := s.NewRows([]string{
					"id",
					"name",
					"email",
					"password",
					"is_admin",
					"created_at",
					"updated_at",
				}).AddRow(
					1,
					"Younes",
					"younes@example.com",
					"hashed-password",
					false,
					now,
					nil,
				)

				s.ExpectQuery("SELECT * FROM users WHERE email=$1").
					WithArgs("younes@example.com").
					WillReturnRows(rows)

				got, err := ps.GetUser(context.Background(), "younes@example.com")

				require.NoError(t, err)
				require.Equal(t, int64(1), got.ID)
				require.Equal(t, "Younes", got.Name)
				require.Equal(t, "younes@example.com", got.Email)

				require.NoError(t, s.ExpectationsWereMet())
			},
		},
		{
			name: "failed getting user",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectQuery("SELECT * FROM users WHERE email=$1").
					WithArgs("younes@example.com").
					WillReturnError(fmt.Errorf("db error"))

				_, err := ps.GetUser(context.Background(), "younes@example.com")

				require.Error(t, err)
				require.NoError(t, s.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, s)
			})
		})
	}
}

func TestListUsers(t *testing.T) {
	values := [][]driver.Value{
		{1, "Younes", "younes@example.com", "hash1", false, nil},
		{2, "Admin", "admin@example.com", "hash2", true, nil},
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				rows := s.NewRows([]string{
					"id",
					"name",
					"email",
					"password",
					"is_admin",
					"updated_at",
				}).AddRows(values...)

				s.ExpectQuery("SELECT * FROM users").
					WillReturnRows(rows)

				users, err := ps.ListUsers(context.Background())

				require.NoError(t, err)
				require.Len(t, users, 2)
				require.Equal(t, int64(1), users[0].ID)
				require.Equal(t, int64(2), users[1].ID)

				require.NoError(t, s.ExpectationsWereMet())
			},
		},
		{
			name: "failed listing users",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectQuery("SELECT * FROM users").
					WillReturnError(fmt.Errorf("db error"))

				_, err := ps.ListUsers(context.Background())

				require.Error(t, err)
				require.NoError(t, s.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, s)
			})
		})
	}
}

func TestUpdateUser(t *testing.T) {
	now := time.Now()

	u := &model.User{
		ID:        1,
		Name:      "Updated Younes",
		Email:     "updated@example.com",
		Password:  "new-hash",
		IsAdmin:   true,
		UpdatedAt: &now,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectExec(`
					UPDATE users
					SET
						name=$1,
						email=$2,
						password=$3,
						is_admin=$4,
						updated_at=$5
					WHERE id=$6
				`).
					WithArgs(
						u.Name,
						u.Email,
						u.Password,
						u.IsAdmin,
						u.UpdatedAt,
						u.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				got, err := ps.UpdateUser(context.Background(), u)

				require.NoError(t, err)
				require.Equal(t, u, got)

				require.NoError(t, s.ExpectationsWereMet())
			},
		},
		{
			name: "failed updating user",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectExec(`
					UPDATE users
					SET
						name=$1,
						email=$2,
						password=$3,
						is_admin=$4,
						updated_at=$5
					WHERE id=$6
				`).
					WithArgs(
						u.Name,
						u.Email,
						u.Password,
						u.IsAdmin,
						u.UpdatedAt,
						u.ID,
					).
					WillReturnError(fmt.Errorf("db error"))

				_, err := ps.UpdateUser(context.Background(), u)

				require.Error(t, err)
				require.NoError(t, s.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, s)
			})
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectExec("DELETE FROM users WHERE id=$1").
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))

				err := ps.DeleteUser(context.Background(), 1)

				require.NoError(t, err)
				require.NoError(t, s.ExpectationsWereMet())
			},
		},
		{
			name: "failed deleting user",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectExec("DELETE FROM users WHERE id=$1").
					WithArgs(int64(1)).
					WillReturnError(fmt.Errorf("db error"))

				err := ps.DeleteUser(context.Background(), 1)

				require.Error(t, err)
				require.NoError(t, s.ExpectationsWereMet())
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
				ps := NewPostgresStorer(db)
				tc.test(t, ps, s)
			})
		})
	}
}