package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestDefaultJWTConfig(t *testing.T) {
	c := DefaultJWTConfig("secret")

	require.Equal(t, "secret", c.SecretKey)
	require.Equal(t, 15*time.Minute, c.AccessTokenExpiry)
	require.Equal(t, 7*24*time.Hour, c.RefreshTokenExpiry)
	require.Equal(t, "go-ecommerce", c.Issuer)
}

func TestGenerateAccessToken(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				token, err := c.GenerateAccessToken(1, "younes@example.com", true)

				require.NoError(t, err)
				require.NotEmpty(t, token)

				claims, err := c.ValidateToken(token, "access")
				require.NoError(t, err)

				require.Equal(t, int64(1), claims.ID)
				require.Equal(t, "younes@example.com", claims.Email)
				require.True(t, claims.IsAdmin)
				require.Equal(t, "access", claims.TokenType)
				require.Equal(t, "go-ecommerce", claims.Issuer)
				require.Equal(t, "younes@example.com", claims.Subject)
				require.NotEmpty(t, claims.RegisteredClaims.ID)
				require.WithinDuration(
					t,
					time.Now().Add(c.AccessTokenExpiry),
					claims.ExpiresAt.Time,
					time.Minute,
				)
			},
		},
		{
			name: "every token gets its own id",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				first, err := c.GenerateAccessToken(1, "younes@example.com", false)
				require.NoError(t, err)

				second, err := c.GenerateAccessToken(1, "younes@example.com", false)
				require.NoError(t, err)

				firstClaims, err := c.ValidateToken(first, "access")
				require.NoError(t, err)

				secondClaims, err := c.ValidateToken(second, "access")
				require.NoError(t, err)

				require.NotEqual(
					t,
					firstClaims.RegisteredClaims.ID,
					secondClaims.RegisteredClaims.ID,
				)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	c := DefaultJWTConfig("secret")

	token, err := c.GenerateRefreshToken(1, "younes@example.com", false)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := c.ValidateToken(token, "refresh")
	require.NoError(t, err)

	require.Equal(t, int64(1), claims.ID)
	require.Equal(t, "younes@example.com", claims.Email)
	require.False(t, claims.IsAdmin)
	require.Equal(t, "refresh", claims.TokenType)
	require.WithinDuration(
		t,
		time.Now().Add(c.RefreshTokenExpiry),
		claims.ExpiresAt.Time,
		time.Minute,
	)
}

func TestValidateToken(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				token, err := c.GenerateAccessToken(1, "younes@example.com", false)
				require.NoError(t, err)

				claims, err := c.ValidateToken(token, "access")

				require.NoError(t, err)
				require.Equal(t, int64(1), claims.ID)
			},
		},
		{
			name: "strips the bearer prefix",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				token, err := c.GenerateAccessToken(1, "younes@example.com", false)
				require.NoError(t, err)

				claims, err := c.ValidateToken("Bearer "+token, "access")

				require.NoError(t, err)
				require.Equal(t, int64(1), claims.ID)
			},
		},
		{
			name: "trims surrounding whitespace",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				token, err := c.GenerateAccessToken(1, "younes@example.com", false)
				require.NoError(t, err)

				claims, err := c.ValidateToken("Bearer  "+token+"  ", "access")

				require.NoError(t, err)
				require.Equal(t, int64(1), claims.ID)
			},
		},
		{
			name: "malformed token",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				_, err := c.ValidateToken("not-a-token", "access")

				require.ErrorIs(t, err, ErrInvalidToken)
			},
		},
		{
			name: "empty token",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				_, err := c.ValidateToken("", "access")

				require.ErrorIs(t, err, ErrInvalidToken)
			},
		},
		{
			name: "expired token",
			test: func(t *testing.T) {
				c := &JWTConfig{
					SecretKey:         "secret",
					AccessTokenExpiry: -time.Hour,
					Issuer:            "go-ecommerce",
				}

				token, err := c.GenerateAccessToken(1, "younes@example.com", false)
				require.NoError(t, err)

				_, err = c.ValidateToken(token, "access")

				require.ErrorIs(t, err, ErrTokenExpired)
			},
		},
		{
			name: "signed with another secret",
			test: func(t *testing.T) {
				signer := DefaultJWTConfig("another-secret")

				token, err := signer.GenerateAccessToken(1, "younes@example.com", false)
				require.NoError(t, err)

				verifier := DefaultJWTConfig("secret")

				_, err = verifier.ValidateToken(token, "access")

				require.ErrorIs(t, err, ErrInvalidToken)
			},
		},
		{
			name: "wrong token type",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				token, err := c.GenerateRefreshToken(1, "younes@example.com", false)
				require.NoError(t, err)

				_, err = c.ValidateToken(token, "access")

				require.ErrorIs(t, err, ErrInvalidToken)
			},
		},
		{
			name: "unsigned token is rejected",
			test: func(t *testing.T) {
				c := DefaultJWTConfig("secret")

				now := time.Now()

				claims := Claims{
					ID:        1,
					Email:     "younes@example.com",
					TokenType: "access",
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    c.Issuer,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
					},
				}

				token, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).
					SignedString(jwt.UnsafeAllowNoneSignatureType)
				require.NoError(t, err)

				_, err = c.ValidateToken(token, "access")

				require.ErrorIs(t, err, ErrInvalidToken)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}
