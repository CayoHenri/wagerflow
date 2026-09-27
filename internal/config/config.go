package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName  string
	HTTPPort string
	Env      string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppName:  getEnv("APP_NAME", "wagerflow"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		Env:      getEnv("APP_ENV", "development"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
