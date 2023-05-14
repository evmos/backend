package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
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

// LoadConfig loads the application configuration from environment variables
// or default values specified in the config.toml file.
func LoadConfig() (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

    // Build the absolute path to the target file
    filePath := filepath.Join(dir, "internal/config/config.toml")

	cfg := &Config{}
	_, err = toml.DecodeFile(filePath, cfg)
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
