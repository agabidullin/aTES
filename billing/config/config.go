package config

import "os"

type Config struct {
	DSN string
}

func New() *Config {
	return &Config{
		DSN: getEnv("DSN", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
