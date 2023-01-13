package config

type SentryOpts struct {
	Enabled bool   `yaml:"enabled"`
	DSN     string `yaml:"dsn"`
}

type SuperAdmin struct {
	Email string `env:"SUPER_ADMIN_EMAIL" validate:"required,email"`
	Name  string `yaml:"name"`
}

type Auth struct {
	FakeAuthEnabled bool   `yaml:"fake_auth_enabled"`
	JWTSecret       string `env:"JWT_SECRET"`
	FakeAuthHeaders string `yaml:"fake_auth_headers"`
}

type Folders struct {
	Achievements string `yaml:"achievements"`
	Events       string `yaml:"events"`
	Users        string `yaml:"users"`
}
