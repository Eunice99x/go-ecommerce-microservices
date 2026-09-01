package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const minSecretKeySize = 32

type Config struct {
	SecretKey string `mapstructure:"SECRET_KEY"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError

		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		_ = notFound
		log.Println("warning: no .env file found, using environment variables")
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.SecretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY is required")
	}

	if len(cfg.SecretKey) < minSecretKeySize {
		return nil, fmt.Errorf(
			"SECRET_KEY must be at least %d characters",
			minSecretKeySize,
		)
	}

	return &cfg, nil
}
