package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/santosjordi/fc_challenge_weather_cloudrun/config"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a mock .env file inside the temporary directory
	tempFile := filepath.Join(tempDir, ".env")
	content := `
SERVER_PORT=:9000
WEATHER_API_KEY=test_key_123
`
	err := os.WriteFile(tempFile, []byte(content), 0644)
	require.NoError(t, err)

	// Temporarily override environment variables
	os.Setenv("SERVER_PORT", "")
	os.Setenv("WEATHER_API_KEY", "")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("WEATHER_API_KEY")
	}()

	// Tell viper to use the test .env file by setting working dir
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tempDir)

	// Call the function under test
	cfg, err := config.LoadConfig()

	// Assert that no error occurred
	require.NoError(t, err)

	// Assert that the loaded config matches the expected values
	require.Equal(t, ":9000", cfg.ServerPort)
	require.Equal(t, "test_key_123", cfg.WeatherAPIKey)
}
