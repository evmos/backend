package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strconv"
)

// Config represents the application configuration.
type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
}

// ServerConfig represents the server configuration.
type ServerConfig struct {
	Port       int
	ListenAddr string
}

// LoggingConfig represents the log configuration.
type LoggingConfig struct {
	Level string
	File  string
}

// LoadConfig loads the application configuration from environment variables or default values specified in the config.toml file.
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	f := "config.toml"
	if _, err := os.Stat(f); err != nil {
		f = "_config/config.toml"
	}

	absPath, _ := filepath.Abs("../backend/internal/config/config.toml")

	_, err := toml.DecodeFile(absPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode default config file: %w", err)
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err == nil && port != 0 {
		cfg.Server.Port = port
	}

	logLevel := os.Getenv("LOGGING_LEVEL")
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = logLevel
	}

	loggingFile := os.Getenv("LOGGING_FILE")
	if cfg.Logging.File == "" {
		cfg.Logging.File = loggingFile
	}

	return cfg, nil
}
