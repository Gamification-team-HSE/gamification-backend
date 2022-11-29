package config

import (
	"gitlab.com/krespix/gamification-api/internal/clients/smtp"
	"gitlab.com/krespix/gamification-api/internal/core/listeners/http"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres"
)

// AppConfig represents the configuration of our application.
type AppConfig struct {
	HTTP       http.Config     `yaml:"http"`
	DB         postgres.Config `yaml:"db"`
	Sentry     SentryOpts      `yaml:"sentry"`
	SuperAdmin SuperAdmin      `yaml:"super_admin"`
	SMTP       smtp.Options    `yaml:"smtp"`
	JWT        JWT
}
