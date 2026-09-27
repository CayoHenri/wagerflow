package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("should load default values", func(t *testing.T) {
		t.Setenv("APP_NAME", "")
		t.Setenv("APP_ENV", "")
		t.Setenv("HTTP_PORT", "")

		cfg, err := Load()

		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, defaultAppName, cfg.App.Name)
		assert.Equal(t, defaultAppEnv, cfg.App.Env)
		assert.Equal(t, defaultHTTPPort, cfg.HTTP.Port)
	})

	t.Run("should load values from environment", func(t *testing.T) {
		t.Setenv("APP_NAME", "wagerflow-test")
		t.Setenv("APP_ENV", "test")
		t.Setenv("HTTP_PORT", "9090")

		cfg, err := Load()

		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "wagerflow-test", cfg.App.Name)
		assert.Equal(t, "test", cfg.App.Env)
		assert.Equal(t, 9090, cfg.HTTP.Port)
	})

	t.Run("should return error when HTTP_PORT is not a number", func(t *testing.T) {
		t.Setenv("HTTP_PORT", "invalid")

		cfg, err := Load()

		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.ErrorContains(t, err, "HTTP_PORT must be a valid integer")
	})

	t.Run("should return error when HTTP_PORT is less than one", func(t *testing.T) {
		t.Setenv("HTTP_PORT", "0")

		cfg, err := Load()

		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.ErrorContains(t, err, "HTTP_PORT must be between 1 and 65535")
	})

	t.Run("should return error when HTTP_PORT is greater than 65535", func(t *testing.T) {
		t.Setenv("HTTP_PORT", "65536")

		cfg, err := Load()

		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.ErrorContains(t, err, "HTTP_PORT must be between 1 and 65535")
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

func TestGetEnvAsInt(t *testing.T) {
	t.Run("should return environment value as integer", func(t *testing.T) {
		t.Setenv("TEST_INT", "123")

		value, err := getEnvAsInt("TEST_INT", 10)

		require.NoError(t, err)
		assert.Equal(t, 123, value)
	})

	t.Run("should return fallback when environment value is empty", func(t *testing.T) {
		t.Setenv("TEST_INT", "")

		value, err := getEnvAsInt("TEST_INT", 10)

		require.NoError(t, err)
		assert.Equal(t, 10, value)
	})

	t.Run("should return error when environment value is not an integer", func(t *testing.T) {
		t.Setenv("TEST_INT", "abc")

		value, err := getEnvAsInt("TEST_INT", 10)

		require.Error(t, err)
		assert.Zero(t, value)
		assert.ErrorContains(t, err, "TEST_INT must be a valid integer")
	})
}
