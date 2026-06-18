// Package config loads runtime configuration from the environment, applying
// sensible defaults so the service runs locally with zero setup.
package config

import "os"

type Config struct {
	// Port is the HTTP listen port.
	Port string
	// DatabaseURL is the SQLite file path.
	DatabaseURL string
}

// Load reads configuration from the environment.
func Load() Config {
	return Config{
		Port:           getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "data/pismo.db"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
