package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var configPath string

// Initialize the config
func Init() {
	configDir, _ := os.UserConfigDir()
	passmanDir := filepath.Join(configDir, "passman")
	configPath = filepath.Join(passmanDir, "config.yaml")

	viper.AddConfigPath(passmanDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Default values (no api_base_url default — must be explicitly set)
	viper.SetDefault("default_vault", "personal")

	if err := viper.ReadInConfig(); err != nil {
		// No config file yet — fine on first run
		fmt.Println("No config file found. Run 'pman config set api_base_url <your-server-url>' to set your backend url.")
	}
}

// Get retrieves a value
func Get(key string) string {
	return viper.GetString(key)
}

// Set sets a value and saves the config
func Set(key, value string) error {
	viper.Set(key, value)

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return viper.WriteConfigAs(configPath)
}
