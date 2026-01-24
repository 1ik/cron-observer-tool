package config

import "time"

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Gmail    GmailConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URI      string        `mapstructure:"uri"`
	Name     string        `mapstructure:"name"`
	Timeout  time.Duration `mapstructure:"timeout"`
	MaxConns int           `mapstructure:"max_conns"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret   string   `mapstructure:"jwt_secret"`
	SuperAdmins []string `mapstructure:"super_admins"` // Comma-separated list of super admin emails
}

// GmailConfig holds Gmail SMTP configuration
type GmailConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}
