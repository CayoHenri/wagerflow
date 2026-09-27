package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("should load default values", func(t *testing.T) {
		t.Setenv("APP_NAME", "")
		t.Setenv("HTTP_PORT", "")
		t.Setenv("APP_ENV", "")

		cfg, err := Load()

		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "wagerflow", cfg.AppName)
		assert.Equal(t, "8080", cfg.HTTPPort)
		assert.Equal(t, "development", cfg.Env)
	})

	t.Run("should load values from environment", func(t *testing.T) {
		t.Setenv("APP_NAME", "wagerflow-test")
		t.Setenv("HTTP_PORT", "9090")
		t.Setenv("APP_ENV", "test")

		cfg, err := Load()

		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "wagerflow-test", cfg.AppName)
		assert.Equal(t, "9090", cfg.HTTPPort)
		assert.Equal(t, "test", cfg.Env)
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("should return environment value when defined", func(t *testing.T) {
		t.Setenv("TEST_VALUE", "custom-value")

		value := getEnv("TEST_VALUE", "fallback")

		assert.Equal(t, "custom-value", value)
	})

	t.Run("should return fallback when environment value is empty", func(t *testing.T) {
		t.Setenv("TEST_VALUE", "")

		value := getEnv("TEST_VALUE", "fallback")

		assert.Equal(t, "fallback", value)
	})
}
