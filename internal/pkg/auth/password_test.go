package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				hashed, err := HashPassword("plain-password")

				require.NoError(t, err)
				require.NotEmpty(t, hashed)
				require.NotEqual(t, "plain-password", hashed)
				require.NoError(t, ComparePassword("plain-password", hashed))
			},
		},
		{
			name: "same password hashes to a different value every time",
			test: func(t *testing.T) {
				first, err := HashPassword("plain-password")
				require.NoError(t, err)

				second, err := HashPassword("plain-password")
				require.NoError(t, err)

				require.NotEqual(t, first, second)
				require.NoError(t, ComparePassword("plain-password", first))
				require.NoError(t, ComparePassword("plain-password", second))
			},
		},
		{
			name: "failed hashing a too long password",
			test: func(t *testing.T) {
				// bcrypt rejects anything over 72 bytes
				_, err := HashPassword(strings.Repeat("a", 73))

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestComparePassword(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				hashed, err := HashPassword("plain-password")
				require.NoError(t, err)

				require.NoError(t, ComparePassword("plain-password", hashed))
			},
		},
		{
			name: "wrong password",
			test: func(t *testing.T) {
				hashed, err := HashPassword("plain-password")
				require.NoError(t, err)

				err = ComparePassword("wrong-password", hashed)

				require.Error(t, err)
				require.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
			},
		},
		{
			name: "hash is not a bcrypt hash",
			test: func(t *testing.T) {
				err := ComparePassword("plain-password", "not-a-hash")

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestIsHashedPassword(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "a bcrypt hash",
			test: func(t *testing.T) {
				hashed, err := HashPassword("plain-password")
				require.NoError(t, err)

				require.True(t, IsHashedPassword(hashed))
			},
		},
		{
			name: "a plain password",
			test: func(t *testing.T) {
				require.False(t, IsHashedPassword("plain-password"))
			},
		},
		{
			name: "an empty password",
			test: func(t *testing.T) {
				require.False(t, IsHashedPassword(""))
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}
