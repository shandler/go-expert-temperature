package config_test

import (
	"os"
	"testing"

	"github.com/shandler/go-expert-temperature/internal/config"
)

func TestConfig(t *testing.T) {
	t.Run("Test Config", func(t *testing.T) {
		os.Setenv("WEATHER_KEY", "test")
		config := config.New()

		if config.WeatherKey != "test" {
			t.Errorf("Expected WeatherKey to be test, got %s", config.WeatherKey)
		}
		os.Unsetenv("WEATHER_KEY")
	})
}
