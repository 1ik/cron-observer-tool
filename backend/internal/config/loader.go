package config

import (
	"github.com/spf13/viper"
)

// Load reads configuration from environment variables and .env file
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults for optional fields
	setDefaults(v)

	// Bind environment variables
	bindEnvVars(v)

	// Load .env file (optional, won't error if not found)
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./backend")
	_ = v.ReadInConfig() // Ignore error if .env doesn't exist

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// setDefaults sets default values for optional configuration fields
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")

	// Database defaults (only for optional fields)
	v.SetDefault("database.timeout", "10s")
	v.SetDefault("database.max_conns", 100)
}

// bindEnvVars binds environment variables to configuration keys
func bindEnvVars(v *viper.Viper) {
	// Server environment variables
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	v.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")

	// Database environment variables (required)
	v.BindEnv("database.uri", "DATABASE_URI")
	v.BindEnv("database.name", "DATABASE_NAME")

	// Database environment variables (optional)
	v.BindEnv("database.timeout", "DATABASE_TIMEOUT")
	v.BindEnv("database.max_conns", "DATABASE_MAX_CONNS")
}
