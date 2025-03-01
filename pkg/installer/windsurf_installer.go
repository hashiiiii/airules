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
	templateDir     string
	localDestDir    string
	globalDestDir   string
	localFileName   string
	globalFileName  string
}

// NewWindsurfInstaller returns a new instance of WindsurfInstaller
func NewWindsurfInstaller() *WindsurfInstaller {
	home, _ := homedir.Dir()
	
	// Get Windsurf template path from vendor
	templateDir := filepath.Join("vendor", "rules-for-ai", "windsurf")
	
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
		localFileName:  "cascade.local.json",
		globalFileName: "cascade.global.json",
	}
}

// InstallLocal installs the local configuration file
func (i *WindsurfInstaller) InstallLocal() error {
	srcPath := filepath.Join(i.templateDir, i.localFileName)
	destPath := filepath.Join(i.localDestDir, i.localFileName)
	
	// Check if directory exists, create if not
	if err := os.MkdirAll(i.localDestDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Copy file
	return copyFile(srcPath, destPath)
}

// InstallGlobal installs the global configuration file
func (i *WindsurfInstaller) InstallGlobal() error {
	srcPath := filepath.Join(i.templateDir, i.globalFileName)
	destPath := filepath.Join(i.globalDestDir, i.globalFileName)
	
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
