// Package config provides configuration management for sub2api.
// It handles loading and validating application settings from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration values.
type Config struct {
	// Server settings
	Host string
	Port int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Subscription settings
	SubURL        string
	RefreshInterval time.Duration
	UserAgent     string

	// Cache settings
	CacheEnabled bool
	CacheTTL     time.Duration

	// Logging
	LogLevel string
}

// Load reads configuration from environment variables and returns a Config.
// It returns an error if any required configuration is missing or invalid.
func Load() (*Config, error) {
	cfg := &Config{
		Host:            getEnvOrDefault("HOST", "0.0.0.0"),
		Port:            getEnvIntOrDefault("PORT", 8080),
		ReadTimeout:     getEnvDurationOrDefault("READ_TIMEOUT", 30*time.Second),
		WriteTimeout:    getEnvDurationOrDefault("WRITE_TIMEOUT", 30*time.Second),
		SubURL:          os.Getenv("SUB_URL"),
		RefreshInterval: getEnvDurationOrDefault("REFRESH_INTERVAL", 12*time.Hour),
		UserAgent:       getEnvOrDefault("USER_AGENT", "sub2api/1.0"),
		CacheEnabled:    getEnvBoolOrDefault("CACHE_ENABLED", true),
		CacheTTL:        getEnvDurationOrDefault("CACHE_TTL", 10*time.Minute),
		LogLevel:        getEnvOrDefault("LOG_LEVEL", "info"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks that required fields are set and values are within acceptable ranges.
func (c *Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535, got %d", c.Port)
	}
	if c.RefreshInterval < time.Minute {
		return fmt.Errorf("REFRESH_INTERVAL must be at least 1 minute")
	}
	if c.LogLevel != "debug" && c.LogLevel != "info" && c.LogLevel != "warn" && c.LogLevel != "error" {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}
	return nil
}

// Addr returns the formatted host:port string for the server listener.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvIntOrDefault(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBoolOrDefault(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func getEnvDurationOrDefault(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
