package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestDBError(t *testing.T) {
	other := errors.New("connection reset")

	tcs := []struct {
		name string
		err  error
		want error
	}{
		{name: "no rows", err: fmt.Errorf("scanning: %w", sql.ErrNoRows), want: model.ErrNotFound},
		{name: "unique violation", err: &pq.Error{Code: pqUniqueViolation}, want: model.ErrAlreadyExists},
		{name: "other pq error", err: &pq.Error{Code: "23503"}, want: nil},
		{name: "unrelated error", err: other, want: other},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := dbError(tc.err)

			if tc.want == nil {
				require.Equal(t, tc.err, got)
				return
			}

			require.ErrorIs(t, got, tc.want)
		})
	}
}
