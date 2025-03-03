package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

// WindsurfInstaller represents an installer for Windsurf configuration files
type WindsurfInstaller struct {
	templateDir    string
	localDestDir   string
	globalDestDir  string
	localFileName  string
	globalFileName string
}

// NewWindsurfInstaller returns a new instance of WindsurfInstaller
func NewWindsurfInstaller() *WindsurfInstaller {
	home, _ := homedir.Dir()

	// Get template directory from environment variable or use default
	templateDir := os.Getenv("AIRULES_TEMPLATE_DIR")
	if templateDir == "" {
		templateDir = filepath.Join("templates", "rules-for-ai", "windsurf")
	}

	// Set destination directories based on OS
	var localDestDir, globalDestDir string

	switch runtime.GOOS {
	case "darwin", "linux":
		// For macOS and Linux
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(home, ".config", "windsurf")
	case "windows":
		// For Windows
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(os.Getenv("APPDATA"), "Windsurf")
	default:
		// For other OS
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(home, ".windsurf")
	}

	return &WindsurfInstaller{
		templateDir:    templateDir,
		localDestDir:   localDestDir,
		globalDestDir:  globalDestDir,
		localFileName:  "local/.windsurfrules", // 実際のファイル名に変更
		globalFileName: "global/.windsurfrules",
	}
}

// InstallLocal installs the local configuration file
func (i *WindsurfInstaller) InstallLocal() error {
	srcPath := filepath.Join(i.templateDir, i.localFileName)

	// ディレクトリ構造を作成
	configDir := filepath.Join(i.localDestDir, ".config", "windsurf")
	destPath := filepath.Join(configDir, "config.json")

	// Check if directory exists, create if not
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy file
	return copyFile(srcPath, destPath)
}

// InstallGlobal installs the global configuration file
func (i *WindsurfInstaller) InstallGlobal() error {
	srcPath := filepath.Join(i.templateDir, i.globalFileName)

	// ディレクトリ構造を作成
	configDir := i.globalDestDir
	destPath := filepath.Join(configDir, "config.json")

	// Check if directory exists, create if not
	if err := os.MkdirAll(i.globalDestDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy file
	return copyFile(srcPath, destPath)
}

// InstallAll installs both configuration files
func (i *WindsurfInstaller) InstallAll() error {
	if err := i.InstallLocal(); err != nil {
		return fmt.Errorf("failed to install local configuration file: %w", err)
	}

	if err := i.InstallGlobal(); err != nil {
		return fmt.Errorf("failed to install global configuration file: %w", err)
	}

	return nil
}
