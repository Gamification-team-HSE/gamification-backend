package config

import (
	"gitlab.com/krespix/gamification-api/internal/clients/s3"
	"gitlab.com/krespix/gamification-api/internal/clients/smtp"
	"gitlab.com/krespix/gamification-api/internal/core/listeners/http"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres"
	"gitlab.com/krespix/gamification-api/internal/services/image"
)

// AppConfig represents the configuration of our application.
type AppConfig struct {
	HTTP         http.Config     `yaml:"http"`
	DB           postgres.Config `yaml:"db"`
	Sentry       SentryOpts      `yaml:"sentry"`
	SuperAdmin   SuperAdmin      `yaml:"super_admin"`
	SMTP         smtp.Options    `yaml:"smtp"`
	Auth         Auth            `yaml:"auth"`
	S3           s3.Options      `yaml:"s3"`
	Buckets      Folders         `yaml:"folders"`
	ImageService image.Options   `yaml:"image"`
}
