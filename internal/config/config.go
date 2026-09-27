package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	defaultAppName  = "wagerflow"
	defaultAppEnv   = "development"
	defaultHTTPPort = 8080
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type HTTPConfig struct {
	Port int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	httpPort, err := getEnvAsInt("HTTP_PORT", defaultHTTPPort)
	if err != nil {
		return nil, err
	}

	if httpPort < 1 || httpPort > 65535 {
		return nil, fmt.Errorf("HTTP_PORT must be between 1 and 65535")
	}

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", defaultAppName),
			Env:  getEnv("APP_ENV", defaultAppEnv),
		},
		HTTP: HTTPConfig{
			Port: httpPort,
		},
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

func getEnvAsInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer: %w", key, err)
	}

	return parsed, nil
}
