package email

import (
	"fmt"

	"github.com/pkg/errors"
)

// Config represents the SMTP configuration for email sending.
// These settings should be provided by the self-hosted instance administrator.
type Config struct {
	SMTPHost     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	SMTPPort     int
	UseTLS       bool
	UseSSL       bool
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.SMTPHost == "" {
		return errors.New("SMTP host is required")
	}
	if c.SMTPPort <= 0 || c.SMTPPort > 65535 {
		return errors.New("SMTP port must be between 1 and 65535")
	}
	if c.FromEmail == "" {
		return errors.New("from email is required")
	}
	return nil
}

// GetServerAddress returns the SMTP server address in the format "host:port".
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.SMTPHost, c.SMTPPort)
}
