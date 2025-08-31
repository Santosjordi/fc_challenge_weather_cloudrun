package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the configuration parameters for the rate limiter service.
// It includes settings for request limits, lockout durations, window size, and Redis connection details.
// Fields:
//   - WeatherAPIToken: API key to the weather service.
//   - ServerPort: Port on which the server listens for incoming requests.
type Config struct {
	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
}

func LoadConfig(envFilePath string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(envFilePath) // Use full path like "/project/.env"
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read .env file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
