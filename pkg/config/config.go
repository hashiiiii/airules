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
	Editors map[string]EditorConfig `toml:"editors"`
}

// EditorConfig represents editor-specific configuration
type EditorConfig struct {
	Local  map[string][]string `toml:"local"`
	Global map[string][]string `toml:"global"`
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Editors: map[string]EditorConfig{
			"windsurf": {
				Local: map[string][]string{
					"default": {"templates/windsurf/local/.windsurfrules"},
				},
				Global: map[string][]string{
					"default": {"templates/windsurf/global/global_rules.md"},
				},
			},
			"cursor": {
				Local: map[string][]string{
					"default": {"templates/cursor/local/project_rules.mdc"},
				},
				Global: map[string][]string{
					"default": {"templates/cursor/global/global_rules.mdc"},
				},
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

	// Ensure the editors map is initialized
	if config.Editors == nil {
		config.Editors = make(map[string]EditorConfig)
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

	// Check if editor exists
	editorConfig, ok := config.Editors[editor]
	if !ok {
		return nil, fmt.Errorf("editor '%s' not found", editor)
	}

	var ruleFiles []string

	// Get rule file names from config based on mode
	switch mode {
	case "local":
		ruleFiles, ok = editorConfig.Local[key]
	case "global":
		ruleFiles, ok = editorConfig.Global[key]
	default:
		return nil, fmt.Errorf("invalid mode '%s'", mode)
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

// GetSupportedEditors returns a list of supported editors
func GetSupportedEditors() []string {
	config, err := LoadConfig()
	if err != nil {
		return []string{"windsurf", "cursor"} // Default editors if config can't be loaded
	}

	editors := make([]string, 0, len(config.Editors))
	for editor := range config.Editors {
		editors = append(editors, editor)
	}
	return editors
}
