package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

// CursorInstaller represents an installer for Cursor configuration files
type CursorInstaller struct {
	templateDir     string
	localDestDir    string
	globalDestDir   string
	localFileName   string
	globalFileName  string
}

// NewCursorInstaller returns a new instance of CursorInstaller
func NewCursorInstaller() *CursorInstaller {
	home, _ := homedir.Dir()
	
	// Get Cursor template path from vendor
	templateDir := filepath.Join("vendor", "rules-for-ai", "cursor")
	
	// Set destination directories based on OS
	var localDestDir, globalDestDir string
	
	switch runtime.GOOS {
	case "darwin":
		// For macOS
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(home, "Library", "Application Support", "Cursor")
	case "linux":
		// For Linux
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(home, ".config", "cursor")
	case "windows":
		// For Windows
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(os.Getenv("APPDATA"), "Cursor")
	default:
		// For other OS
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(home, ".cursor")
	}
	
	return &CursorInstaller{
		templateDir:    templateDir,
		localDestDir:   localDestDir,
		globalDestDir:  globalDestDir,
		localFileName:  "prompt_library.local.json",
		globalFileName: "prompt_library.global.json",
	}
}

// InstallLocal installs the local configuration file
func (i *CursorInstaller) InstallLocal() error {
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
func (i *CursorInstaller) InstallGlobal() error {
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
func (i *CursorInstaller) InstallAll() error {
	if err := i.InstallLocal(); err != nil {
		return fmt.Errorf("failed to install local configuration file: %w", err)
	}
	
	if err := i.InstallGlobal(); err != nil {
		return fmt.Errorf("failed to install global configuration file: %w", err)
	}
	
	return nil
}
