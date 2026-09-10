package config

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

const minSecretKeySize = 32

type Config struct {
	AppEnv     string `mapstructure:"APP_ENV"`
	ServerPort string `mapstructure:"SERVER_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	SecretKey string `mapstructure:"SECRET_KEY"`
}

var defaults = map[string]string{
	"APP_ENV":     "development",
	"SERVER_PORT": "3000",
	"DB_HOST":     "localhost",
	"DB_PORT":     "5433",
	"DB_USER":     "postgres",
	"DB_PASSWORD": "postgres",
	"DB_NAME":     "ecomm",
	"DB_SSLMODE":  "disable",
	"SECRET_KEY":  "",
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	// otherwise DB_NAME="" is treated as unset and silently falls back to the default
	v.AllowEmptyEnv(true)

	for key, value := range defaults {
		v.SetDefault(key, value)

		// Unmarshal only sees env vars that are bound, not AutomaticEnv ones
		if err := v.BindEnv(key); err != nil {
			return nil, fmt.Errorf("failed to bind %s: %w", key, err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.SecretKey == "" {
		return fmt.Errorf("SECRET_KEY is required")
	}

	if len(c.SecretKey) < minSecretKeySize {
		return fmt.Errorf("SECRET_KEY must be at least %d characters", minSecretKeySize)
	}

	if c.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	if c.DBHost == "" || c.DBUser == "" || c.DBName == "" {
		return fmt.Errorf("DB_HOST, DB_USER and DB_NAME are required")
	}

	return nil
}

// url.URL escapes passwords containing ':' or '@', which would corrupt the DSN
func (c *Config) DSN() string {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.DBUser, c.DBPassword),
		Host:     c.DBHost + ":" + c.DBPort,
		Path:     c.DBName,
		RawQuery: url.Values{"sslmode": {c.DBSSLMode}}.Encode(),
	}

	return dsn.String()
}

func (c *Config) Addr() string {
	return ":" + c.ServerPort
}
