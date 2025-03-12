package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

// Config represents the application configuration
type Config struct {
	Windsurf WindsurfConfig `toml:"windsurf"`
}

// WindsurfConfig represents Windsurf-specific configuration
type WindsurfConfig struct {
	Keys map[string][]string `toml:"keys"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Windsurf: WindsurfConfig{
			Keys: map[string][]string{
				"default": {".windsurfrules"},
			},
		},
	}
}

// GetConfigDir returns the base configuration directory
func GetConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "airules")
	return configDir, nil
}

// EnsureConfigDir creates the configuration directory if it doesn't exist
func EnsureConfigDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

// LoadConfig loads the configuration from file
func LoadConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(configDir, "config.toml")

	// If config file doesn't exist, create default
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := SaveConfig(config); err != nil {
			return nil, err
		}
		return config, nil
	}

	// Read and parse config file
	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	configDir, err := EnsureConfigDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "config.toml")

	// Create or truncate config file
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write config to file
	encoder := toml.NewEncoder(f)
	return encoder.Encode(config)
}

// GetRuleFilePaths returns the paths to rule files by key
func GetRuleFilePaths(key string) ([]string, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Get rule file names from config
	ruleFiles, ok := config.Windsurf.Keys[key]
	if !ok {
		return nil, fmt.Errorf("rule key '%s' not found", key)
	}

	// Get config directory
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	// Convert relative paths to absolute paths
	var absolutePaths []string
	for _, file := range ruleFiles {
		absolutePaths = append(absolutePaths, filepath.Join(configDir, file))
	}

	return absolutePaths, nil
}
