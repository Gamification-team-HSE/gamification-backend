package config

type SentryOpts struct {
	Enabled bool   `yaml:"enabled"`
	DSN     string `yaml:"dsn"`
}
