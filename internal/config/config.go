package config

import "os"

type Config struct {
	WeatherKey string
}

func New() *Config {
	return &Config{
		WeatherKey: os.Getenv("WEATHER_KEY"),
	}
}
