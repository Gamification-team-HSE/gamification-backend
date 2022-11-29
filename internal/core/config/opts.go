package config

type SentryOpts struct {
	Enabled bool   `yaml:"enabled"`
	DSN     string `yaml:"dsn"`
}

type SuperAdmin struct {
	Email string `env:"SUPER_ADMIN_EMAIL" validate:"required,email"`
	Name  string `yaml:"name"`
}

type JWT struct {
	Secret string `env:"JWT_SECRET"`
}
