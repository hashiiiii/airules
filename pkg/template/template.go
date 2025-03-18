package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashiiiii/airules/pkg/config"
)

// ListTemplates lists all available templates for the specified editor and mode
func ListTemplates(editor, mode string) ([]string, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	// Get the appropriate map based on editor and mode
	var rulesMap map[string][]string
	editorConfig, ok := cfg.Editors[editor]
	if !ok {
		return nil, fmt.Errorf("editor '%s' not found", editor)
	}

	if mode == "local" {
		rulesMap = editorConfig.Local
	} else if mode == "global" {
		rulesMap = editorConfig.Global
	} else {
		return nil, fmt.Errorf("invalid mode '%s'", mode)
	}

	// Extract template keys
	keys := make([]string, 0, len(rulesMap))
	for key := range rulesMap {
		keys = append(keys, key)
	}

	return keys, nil
}

// ShowTemplate returns the content of a template
func ShowTemplate(editor, mode, key string) (string, error) {
	// Get rule file paths
	rulePaths, err := config.GetRuleFilePaths(editor, mode, key)
	if err != nil {
		return "", fmt.Errorf("error getting rule file paths: %w", err)
	}

	// Combine content from all rule files
	var content strings.Builder
	for i, path := range rulePaths {
		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return "", fmt.Errorf("template file '%s' not found", path)
		}

		// Read file content
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("error reading template file '%s': %w", path, err)
		}

		// Add separator between files
		if i > 0 {
			content.WriteString("\n\n")
		}

		// Add file path as comment
		content.WriteString(fmt.Sprintf("// From %s\n", filepath.Base(path)))
		content.Write(fileContent)
	}

	return content.String(), nil
}

// ImportTemplate imports a template from a file
func ImportTemplate(editor, mode, key, filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file '%s' not found", filePath)
	}

	// Get config directory
	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("error getting config directory: %w", err)
	}

	// Create templates directory if it doesn't exist
	templatesDir := filepath.Join(configDir, "templates", editor, mode)
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("error creating templates directory: %w", err)
	}

	// Determine file extension based on editor and mode
	var ext string
	switch editor {
	case "windsurf":
		if mode == "local" {
			ext = ".windsurfrules"
		} else {
			ext = ".md"
		}
	case "cursor":
		ext = ".mdc"
	default:
		return fmt.Errorf("unsupported editor: %s", editor)
	}

	// Create destination file path
	destFileName := key + ext
	destPath := filepath.Join(templatesDir, destFileName)

	// Copy file
	srcContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading source file: %w", err)
	}

	if err := os.WriteFile(destPath, srcContent, 0644); err != nil {
		return fmt.Errorf("error writing destination file: %w", err)
	}

	// Update configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Get the appropriate map based on editor and mode
	var rulesMap map[string][]string
	editorConfig, ok := cfg.Editors[editor]
	if !ok {
		editorConfig = config.EditorConfig{
			Local:  make(map[string][]string),
			Global: make(map[string][]string),
		}
		cfg.Editors[editor] = editorConfig
	}

	if mode == "local" {
		rulesMap = editorConfig.Local
		if rulesMap == nil {
			rulesMap = make(map[string][]string)
			editorConfig.Local = rulesMap
		}
	} else {
		rulesMap = editorConfig.Global
		if rulesMap == nil {
			rulesMap = make(map[string][]string)
			editorConfig.Global = rulesMap
		}
	}

	// Add file to key
	relativePath := filepath.Join("templates", editor, mode, destFileName)
	rulesMap[key] = []string{relativePath}

	// Save configuration
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	return nil
}

// ExportTemplate exports a template to a file
func ExportTemplate(editor, mode, key, outputPath string) error {
	// Get template content
	content, err := ShowTemplate(editor, mode, key)
	if err != nil {
		return fmt.Errorf("error getting template content: %w", err)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Write content to output file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	return nil
}
