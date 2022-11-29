package smtp

import (
	"fmt"
	"net/smtp"
)

type Client interface {
	Send(to string, msg []byte) error
}

type Options struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Email    string `yaml:"email"`
	Password string `env:"SMTP_PASSWORD"`
}

type client struct {
	auth smtp.Auth
	addr string
	opts Options
}

func (c *client) Send(to string, msg []byte) error {
	return smtp.SendMail(c.addr, c.auth, c.opts.Email, []string{to}, msg)
}

func New(opts Options) Client {
	return &client{
		opts: opts,
		addr: fmt.Sprintf("%s:%s", opts.Host, opts.Port),
		auth: smtp.PlainAuth("", opts.Email, opts.Password, opts.Host),
	}
}
