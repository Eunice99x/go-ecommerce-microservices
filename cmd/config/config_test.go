package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success with defaults",
			test: func(t *testing.T) {
				t.Setenv("SECRET_KEY", "a-secret-key-that-is-long-enough")

				cfg, err := LoadConfig()

				require.NoError(t, err)
				require.Equal(t, "development", cfg.AppEnv)
				require.Equal(t, "3000", cfg.ServerPort)
				require.Equal(t, "localhost", cfg.DBHost)
				require.Equal(t, "ecomm", cfg.DBName)
			},
		},
		{
			name: "environment overrides defaults",
			test: func(t *testing.T) {
				t.Setenv("SECRET_KEY", "a-secret-key-that-is-long-enough")
				t.Setenv("APP_ENV", "production")
				t.Setenv("SERVER_PORT", "8080")
				t.Setenv("DB_HOST", "db.internal")
				t.Setenv("DB_NAME", "shop")

				cfg, err := LoadConfig()

				require.NoError(t, err)
				require.Equal(t, "production", cfg.AppEnv)
				require.Equal(t, "8080", cfg.ServerPort)
				require.Equal(t, "db.internal", cfg.DBHost)
				require.Equal(t, "shop", cfg.DBName)
			},
		},
		{
			name: "missing secret key",
			test: func(t *testing.T) {
				t.Setenv("SECRET_KEY", "")

				_, err := LoadConfig()

				require.Error(t, err)
				require.Contains(t, err.Error(), "SECRET_KEY is required")
			},
		},
		{
			name: "secret key too short",
			test: func(t *testing.T) {
				t.Setenv("SECRET_KEY", "too-short")

				_, err := LoadConfig()

				require.Error(t, err)
				require.Contains(t, err.Error(), "at least 32 characters")
			},
		},
		{
			name: "missing db name",
			test: func(t *testing.T) {
				t.Setenv("SECRET_KEY", "a-secret-key-that-is-long-enough")
				t.Setenv("DB_NAME", "")

				_, err := LoadConfig()

				require.Error(t, err)
				require.Contains(t, err.Error(), "DB_NAME")
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestDSN(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				cfg := &Config{
					DBHost:     "localhost",
					DBPort:     "5433",
					DBUser:     "postgres",
					DBPassword: "postgres",
					DBName:     "ecomm",
					DBSSLMode:  "disable",
				}

				require.Equal(
					t,
					"postgres://postgres:postgres@localhost:5433/ecomm?sslmode=disable",
					cfg.DSN(),
				)
			},
		},
		{
			name: "escapes a password with reserved characters",
			test: func(t *testing.T) {
				cfg := &Config{
					DBHost:     "localhost",
					DBPort:     "5433",
					DBUser:     "postgres",
					DBPassword: "p@ss:w/rd",
					DBName:     "ecomm",
					DBSSLMode:  "require",
				}

				require.Equal(
					t,
					"postgres://postgres:p%40ss%3Aw%2Frd@localhost:5433/ecomm?sslmode=require",
					cfg.DSN(),
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

func TestAddr(t *testing.T) {
	cfg := &Config{ServerPort: "3000"}

	require.Equal(t, ":3000", cfg.Addr())
}
