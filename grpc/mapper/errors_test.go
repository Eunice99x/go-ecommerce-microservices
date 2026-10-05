package mapper

import (
	"errors"
	"fmt"
	"testing"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorRoundTrip(t *testing.T) {
	for _, e := range errorCodes {
		t.Run(e.code.String(), func(t *testing.T) {
			sent := fmt.Errorf("error getting order: %w", e.err)

			st := ErrorToStatus("/pb.Ecomm/Test", sent)
			require.Equal(t, e.code, status.Code(st))

			got := ErrorFromStatus(st)
			require.ErrorIs(t, got, e.err)
			require.Equal(t, sent.Error(), got.Error())
		})
	}
}

func TestErrorToStatusHidesInternalErrors(t *testing.T) {
	st := ErrorToStatus("/pb.Ecomm/Test", errors.New(`pq: relation "orders" does not exist`))

	require.Equal(t, codes.Internal, status.Code(st))
	require.NotContains(t, st.Error(), "pq:")
}

func TestErrorFromStatusPassesThroughUnknownCodes(t *testing.T) {
	st := status.Error(codes.Unavailable, "connection refused")

	got := ErrorFromStatus(st)

	require.Equal(t, codes.Unavailable, status.Code(got))
	require.NotErrorIs(t, got, model.ErrNotFound)
}
