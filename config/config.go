package config

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port int
}

type EndpointConfig struct {
	Enabled       bool     `yaml:"enabled"`
	EventType     string   `yaml:"event_type"`
	TrackChanges  bool     `yaml:"track_changes"`
	TrackedFields []string `yaml:"tracked_fields"`
}

type ResourceConfig struct {
	Create EndpointConfig `yaml:"create"`
	Update EndpointConfig `yaml:"update"`
	Delete EndpointConfig `yaml:"delete"`
}

type EventTrackingConfig struct {
	Enabled   bool                      `yaml:"enabled"`
	Endpoints map[string]ResourceConfig `yaml:"endpoints"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      auth.JWTConfig
	Redis    struct {
		URL string `yaml:"url"`
	} `yaml:"redis"`
	EventTracking EventTrackingConfig `yaml:"event_tracking"`
}

func LoadConfig() (*Config, error) {
	dbPort, _ := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))

	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "postgres"),
			Port:     dbPort,
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			Name:     getEnvOrDefault("DB_NAME", "aiclinic"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		JWT: auth.JWTConfig{
			Secret:        getEnvOrDefault("JWT_SECRET", "your-256-bit-secret"),
			RefreshSecret: getEnvOrDefault("JWT_REFRESH_SECRET", "your-refresh-secret"),
			ExpiryHours:   24,
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
