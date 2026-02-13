package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Initialize the config
func Init() {
	configDir, _ := os.UserConfigDir()

	// Add config path and set defaults
	viper.AddConfigPath(filepath.Join(configDir, "passman"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Default Values
	viper.SetDefault("api_base_url", "https://shubserver.onrender.com")
	viper.SetDefault("default_vault", "personal")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found, using defaults")
	}
}

// Get retreives a value
func Get(key string) string {
	return viper.GetString(key)
}

// Set sets a value and saves the config
func Set(key, value string) error {
	viper.Set(key, value)
	return viper.WriteConfigAs(filepath.Join(viper.ConfigFileUsed()))
}
