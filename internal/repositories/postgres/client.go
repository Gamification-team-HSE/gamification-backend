package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // imports the postgres driver
	"gitlab.com/krespix/gamification-api/internal/core/errors"
)

const (
	// ErrConnect is returned when we cannot connect to the database.
	ErrConnect = errors.Error("failed to connect to postgres db")
	// ErrClose is returned when we cannot close the database.
	ErrClose = errors.Error("failed to close postgres db connection")
)

// Config represents the configuration for our postgres database.
type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
	Password string `env:"DB_PASSWORD"`
}

// Client provides an implementation for connecting to a postgres database.
type Client struct {
	cfg Config
	db  *sqlx.DB
}

// New instantiates an instance of the Client.
func New(cfg Config) (*Client, error) {
	return &Client{
		cfg: cfg,
	}, nil
}

// Connect connects to the database.
func (c *Client) Connect(_ context.Context) error {
	db, err := sqlx.Connect("postgres", c.createDSN())
	if err != nil {
		return ErrConnect.Wrap(err)
	}
	c.db = db
	return nil
}

// Close closes the database connection.
func (c *Client) Close(_ context.Context) error {
	if err := c.db.Close(); err != nil {
		return ErrClose.Wrap(err)
	}

	return nil
}

// GetDB returns the underlying database connection.
func (c *Client) GetDB() *sql.DB {
	return c.db.DB
}

func (c *Client) GetDBx() *sqlx.DB {
	return c.db
}

func (c *Client) createDSN() string {
	return fmt.Sprintf("user = %s host = %s port = %s password = %s dbname = %s sslmode = %s",
		c.cfg.User, c.cfg.Host, c.cfg.Port, c.cfg.Password, c.cfg.DBName, c.cfg.SSLMode)
}
