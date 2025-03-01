package installer

import (
	"os"
	"path/filepath"
	"testing"
)

// Test for local installation of Cursor installer
func TestCursorInstaller_InstallLocal(t *testing.T) {
	// Setup test environment
	_, templateDir, destDir := setupInstallerTest(t)

	// Create template file
	templateFilePath := filepath.Join(templateDir, "prompt_library.local.json")
	templateContent := []byte(`{"test": "cursor-config"}`)
	err := os.WriteFile(templateFilePath, templateContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Create installer
	installer := NewCursorInstaller()

	// Override installation directories for testing
	installer.templateDir = templateDir
	installer.localDestDir = destDir

	// Test local file installation
	err = installer.InstallLocal()
	if err != nil {
		t.Fatalf("Failed to install local file: %v", err)
	}

	// Check if installed file exists
	destFile := filepath.Join(destDir, installer.localFileName)
	exists := fileExists(destFile)
	if !exists {
		t.Errorf("File was not installed: %s", destFile)
	}

	// Check file content
	installedContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("Failed to read installed file: %v", err)
	}

	if string(installedContent) != string(templateContent) {
		t.Errorf("Installed file content does not match.\nExpected: %s\nActual: %s",
			string(templateContent), string(installedContent))
	}
}

// Test for global installation of Cursor installer
func TestCursorInstaller_InstallGlobal(t *testing.T) {
	// Setup test environment
	_, templateDir, destDir := setupInstallerTest(t)

	// Create template file
	templateFilePath := filepath.Join(templateDir, "prompt_library.global.json")
	templateContent := []byte(`{"test": "cursor-global-config"}`)
	err := os.WriteFile(templateFilePath, templateContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Create installer
	installer := NewCursorInstaller()

	// Override installation directories for testing
	installer.templateDir = templateDir
	installer.globalDestDir = destDir

	// Test global file installation
	err = installer.InstallGlobal()
	if err != nil {
		t.Fatalf("Failed to install global file: %v", err)
	}

	// Check if installed file exists
	destFile := filepath.Join(destDir, installer.globalFileName)
	exists := fileExists(destFile)
	if !exists {
		t.Errorf("File was not installed: %s", destFile)
	}

	// Check file content
	installedContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("Failed to read installed file: %v", err)
	}

	if string(installedContent) != string(templateContent) {
		t.Errorf("Installed file content does not match.\nExpected: %s\nActual: %s",
			string(templateContent), string(installedContent))
	}
}

// Test for installing all files
func TestCursorInstaller_InstallAll(t *testing.T) {
	// Setup test environment
	_, templateDir, destDir := setupInstallerTest(t)

	// Create template files
	localTemplate := filepath.Join(templateDir, "prompt_library.local.json")
	localContent := []byte(`{"test": "local-config"}`)
	err := os.WriteFile(localTemplate, localContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create local template file: %v", err)
	}

	globalTemplate := filepath.Join(templateDir, "prompt_library.global.json")
	globalContent := []byte(`{"test": "global-config"}`)
	err = os.WriteFile(globalTemplate, globalContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create global template file: %v", err)
	}

	// Create installer
	installer := NewCursorInstaller()
	installer.templateDir = templateDir
	installer.localDestDir = destDir
	installer.globalDestDir = destDir

	// Install all files
	err = installer.InstallAll()
	if err != nil {
		t.Fatalf("Failed to install all files: %v", err)
	}

	// Check local file
	localFile := filepath.Join(destDir, installer.localFileName)
	if !fileExists(localFile) {
		t.Errorf("Local file was not installed: %s", localFile)
	}

	// Check local file content
	localInstalledContent, err := os.ReadFile(localFile)
	if err != nil {
		t.Fatalf("Failed to read installed local file: %v", err)
	}

	if string(localInstalledContent) != string(localContent) {
		t.Errorf("Installed local file content does not match.\nExpected: %s\nActual: %s",
			string(localContent), string(localInstalledContent))
	}

	// Check global file
	globalFile := filepath.Join(destDir, installer.globalFileName)
	if !fileExists(globalFile) {
		t.Errorf("Global file was not installed: %s", globalFile)
	}

	// Check global file content
	globalInstalledContent, err := os.ReadFile(globalFile)
	if err != nil {
		t.Fatalf("Failed to read installed global file: %v", err)
	}

	if string(globalInstalledContent) != string(globalContent) {
		t.Errorf("Installed global file content does not match.\nExpected: %s\nActual: %s",
			string(globalContent), string(globalInstalledContent))
	}
}
