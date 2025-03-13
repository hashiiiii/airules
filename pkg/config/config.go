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
	Cursor   CursorConfig   `toml:"cursor"`
}

// WindsurfConfig represents Windsurf-specific configuration
type WindsurfConfig struct {
	Local  map[string][]string `toml:"local"`
	Global map[string][]string `toml:"global"`
}

// CursorConfig represents Cursor-specific configuration
type CursorConfig struct {
	Local  map[string][]string `toml:"local"`
	Global map[string][]string `toml:"global"`
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Windsurf: WindsurfConfig{
			Local: map[string][]string{
				"default": {"templates/windsurf/local/.windsurfrules"},
			},
			Global: map[string][]string{
				"default": {"templates/windsurf/global/global_rules.md"},
			},
		},
		Cursor: CursorConfig{
			Local: map[string][]string{
				"default": {"templates/cursor/local/project_rules.mdc"},
			},
			Global: map[string][]string{
				"default": {"templates/cursor/global/global_rules.mdc"},
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
		return "", fmt.Errorf("Failed to create config directory: %w", err)
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
		config := GetDefaultConfig()
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

// GetRuleFilePaths returns the paths to rule files by editor, mode and key
func GetRuleFilePaths(editor, mode, key string) ([]string, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	var ruleFiles []string
	var ok bool

	// Get rule file names from config based on editor and mode
	switch editor {
	case "windsurf":
		switch mode {
		case "local":
			ruleFiles, ok = config.Windsurf.Local[key]
		case "global":
			ruleFiles, ok = config.Windsurf.Global[key]
		default:
			return nil, fmt.Errorf("invalid mode '%s' for windsurf", mode)
		}
	case "cursor":
		switch mode {
		case "local":
			ruleFiles, ok = config.Cursor.Local[key]
		case "global":
			ruleFiles, ok = config.Cursor.Global[key]
		default:
			return nil, fmt.Errorf("invalid mode '%s' for cursor", mode)
		}
	default:
		return nil, fmt.Errorf("invalid editor '%s'", editor)
	}

	if !ok {
		return nil, fmt.Errorf("rule key '%s' not found for %s %s", key, editor, mode)
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
