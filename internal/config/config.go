package config

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"

	"gitlab.com/krespix/gamification-api/internal/core/config"
	"gitlab.com/krespix/gamification-api/internal/core/errors"
	"gitlab.com/krespix/gamification-api/internal/core/logging"
	"go.uber.org/zap"
)

const (
	// ErrValidation is returned when the configuration is invalid.
	ErrValidation = errors.Error("invalid configuration")
	// ErrEnvVars is returned when the environment variables are invalid.
	ErrEnvVars = errors.Error("failed parsing env vars")
	// ErrRead is returned when the configuration file cannot be read.
	ErrRead = errors.Error("failed to read file")
	// ErrUnmarshal is returned when the configuration file cannot be unmarshalled.
	ErrUnmarshal = errors.Error("failed to unmarshal file")
)

// Config represents the configuration of our application.
type Config struct {
	config.AppConfig `yaml:",inline"`
}

// Load loads the configuration from the config/config.yaml file.
func Load(ctx context.Context, cfgPath string) (*Config, error) {
	cfg := &Config{}

	if err := loadFromFiles(ctx, cfg, cfgPath); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return nil, ErrEnvVars.Wrap(err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, ErrValidation.Wrap(err)
	}

	return cfg, nil
}

func loadFromFiles(ctx context.Context, cfg any, path string) error {
	if err := loadYaml(ctx, path, cfg); err != nil {
		return err
	}
	return nil
}

func loadYaml(ctx context.Context, filename string, cfg any) error {
	logging.From(ctx).Info("Loading configuration", zap.String("path", filename))

	data, err := os.ReadFile(filename)
	if err != nil {
		return ErrRead.Wrap(err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return ErrUnmarshal.Wrap(err)
	}

	return nil
}
