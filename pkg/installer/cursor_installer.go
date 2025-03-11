package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// CursorLocalRulesDir is the directory for local Cursor rules
	CursorLocalRulesDir = ".cursor/rules"

	// CursorLocalRulesFile is the filename for local Cursor rules
	CursorLocalRulesFile = "project_rules.mdc"

	// CursorGlobalRulesFile is the filename for global Cursor rules
	CursorGlobalRulesFile = ".cursorrules"
)

// CursorInstaller implements the Installer interface for Cursor editor
type CursorInstaller struct {
	templateDir string
	language    Language
}

// NewCursorInstaller creates a new CursorInstaller
func NewCursorInstaller(language Language) (*CursorInstaller, error) {
	// Get the template directory from the environment variable
	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		// Default to "templates" in the current directory
		templateDir = "templates"
	}

	// Resolve the template directory to an absolute path
	templateDir, err := filepath.Abs(filepath.Join(templateDir, "cursor"))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template directory: %w", err)
	}

	// Check if the template directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("template directory does not exist: %s", templateDir)
	}

	return &CursorInstaller{
		templateDir: templateDir,
		language:    language,
	}, nil
}

// Install installs Cursor rules based on the installation type
func (i *CursorInstaller) Install(installType InstallType) error {
	switch installType {
	case Local:
		return i.installLocal()
	case Global:
		return i.installGlobal()
	case All:
		if err := i.installLocal(); err != nil {
			return err
		}
		return i.installGlobal()
	default:
		return fmt.Errorf("unknown installation type: %v", installType)
	}
}

// installLocal installs local Cursor rules
func (i *CursorInstaller) installLocal() error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Create the destination directory
	destDir := filepath.Join(cwd, CursorLocalRulesDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Get the source and destination paths
	srcPath, destPath, destDir, err := GetInstallPath(
		Local,
		i.templateDir,
		destDir,
		"", // Not used for local installation
		CursorLocalRulesFile,
		"", // Not used for local installation
		i.language,
	)
	if err != nil {
		return err
	}

	// Check if the destination file already exists
	if _, err := os.Stat(destPath); err == nil {
		// Create a backup of the existing file
		backupPath := destPath + ".backup"
		if err := CopyFile(destPath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup of existing file: %w", err)
		}
		fmt.Printf("Created backup of existing file at %s\n", backupPath)
	}

	// Copy the template file to the destination
	if err := CopyFile(srcPath, destPath); err != nil {
		return fmt.Errorf("failed to copy template file: %w", err)
	}

	fmt.Printf("Installed local Cursor rules to %s\n", destPath)
	return nil
}

// installGlobal installs global Cursor rules
func (i *CursorInstaller) installGlobal() error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// For global installation, we'll create the .cursorrules file in the current directory
	destPath := filepath.Join(cwd, CursorGlobalRulesFile)

	// Get the source path
	srcPath := filepath.Join(i.templateDir, "global", CursorGlobalRulesFile)
	if i.language == Japanese {
		srcPath = filepath.Join(i.templateDir, "global", CursorGlobalRulesFile+"_JA")
	}

	// Check if the source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", srcPath)
	}

	// Check if the destination file already exists
	if _, err := os.Stat(destPath); err == nil {
		// Create a backup of the existing file
		backupPath := destPath + ".backup"
		if err := CopyFile(destPath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup of existing file: %w", err)
		}
		fmt.Printf("Created backup of existing file at %s\n", backupPath)
	}

	// Copy the template file to the destination
	if err := CopyFile(srcPath, destPath); err != nil {
		return fmt.Errorf("failed to copy template file: %w", err)
	}

	fmt.Printf("Installed global Cursor rules to %s\n", destPath)
	return nil
}
