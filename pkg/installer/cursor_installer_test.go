package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewCursorInstaller tests the NewCursorInstaller function
func TestNewCursorInstaller(t *testing.T) {
	// Set up a temporary template directory
	tempDir := t.TempDir()
	cursorDir := filepath.Join(tempDir, "cursor")
	err := os.MkdirAll(filepath.Join(cursorDir, "local"), 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Join(cursorDir, "global"), 0755)
	assert.NoError(t, err)

	// Create test template files
	localFile := filepath.Join(cursorDir, "local", CursorLocalRulesFile)
	err = os.WriteFile(localFile, []byte("# Test local rules"), 0644)
	assert.NoError(t, err)

	globalFile := filepath.Join(cursorDir, "global", CursorGlobalRulesFile)
	err = os.WriteFile(globalFile, []byte("# Test global rules"), 0644)
	assert.NoError(t, err)

	// Set the TEMPLATE_DIR environment variable
	oldTemplateDir := os.Getenv("TEMPLATE_DIR")
	defer os.Setenv("TEMPLATE_DIR", oldTemplateDir)
	os.Setenv("TEMPLATE_DIR", tempDir)

	// Call the function being tested
	installer, err := NewCursorInstaller(English)

	// Assert that there was no error
	assert.NoError(t, err)
	assert.NotNil(t, installer)
	assert.Equal(t, cursorDir, installer.templateDir)
	assert.Equal(t, English, installer.language)
}

// TestNewCursorInstallerWithJapanese tests the NewCursorInstaller function with Japanese language
func TestNewCursorInstallerWithJapanese(t *testing.T) {
	// Set up a temporary template directory
	tempDir := t.TempDir()
	cursorDir := filepath.Join(tempDir, "cursor")
	err := os.MkdirAll(filepath.Join(cursorDir, "local"), 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Join(cursorDir, "global"), 0755)
	assert.NoError(t, err)

	// Create test template files
	localFile := filepath.Join(cursorDir, "local", CursorLocalRulesFile+"_JA")
	err = os.WriteFile(localFile, []byte("# テストローカルルール"), 0644)
	assert.NoError(t, err)

	globalFile := filepath.Join(cursorDir, "global", CursorGlobalRulesFile+"_JA")
	err = os.WriteFile(globalFile, []byte("# テストグローバルルール"), 0644)
	assert.NoError(t, err)

	// Set the TEMPLATE_DIR environment variable
	oldTemplateDir := os.Getenv("TEMPLATE_DIR")
	defer os.Setenv("TEMPLATE_DIR", oldTemplateDir)
	os.Setenv("TEMPLATE_DIR", tempDir)

	// Call the function being tested
	installer, err := NewCursorInstaller(Japanese)

	// Assert that there was no error
	assert.NoError(t, err)
	assert.NotNil(t, installer)
	assert.Equal(t, cursorDir, installer.templateDir)
	assert.Equal(t, Japanese, installer.language)
}

// TestCursorInstallerInstallLocal tests the Install function with Local installation type
func TestCursorInstallerInstallLocal(t *testing.T) {
	// Set up a temporary template directory
	tempDir := t.TempDir()
	cursorDir := filepath.Join(tempDir, "cursor")
	err := os.MkdirAll(filepath.Join(cursorDir, "local"), 0755)
	assert.NoError(t, err)

	// Create test template file
	localFile := filepath.Join(cursorDir, "local", CursorLocalRulesFile)
	err = os.WriteFile(localFile, []byte("# Test local rules"), 0644)
	assert.NoError(t, err)

	// Set the TEMPLATE_DIR environment variable
	oldTemplateDir := os.Getenv("TEMPLATE_DIR")
	defer os.Setenv("TEMPLATE_DIR", oldTemplateDir)
	os.Setenv("TEMPLATE_DIR", tempDir)

	// Create a temporary working directory
	workDir := t.TempDir()
	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(oldWd)
	err = os.Chdir(workDir)
	assert.NoError(t, err)

	// Create the installer
	installer, err := NewCursorInstaller(English)
	assert.NoError(t, err)

	// Call the function being tested
	err = installer.Install(Local)
	assert.NoError(t, err)

	// Check that the file was installed
	destPath := filepath.Join(workDir, CursorLocalRulesDir, CursorLocalRulesFile)
	assert.FileExists(t, destPath)

	// Check the content of the installed file
	content, err := os.ReadFile(destPath)
	assert.NoError(t, err)
	assert.Equal(t, "# Test local rules", string(content))
}

// TestCursorInstallerInstallGlobal tests the Install function with Global installation type
func TestCursorInstallerInstallGlobal(t *testing.T) {
	// Set up a temporary template directory
	tempDir := t.TempDir()
	cursorDir := filepath.Join(tempDir, "cursor")
	err := os.MkdirAll(filepath.Join(cursorDir, "global"), 0755)
	assert.NoError(t, err)

	// Create test template file
	globalFile := filepath.Join(cursorDir, "global", CursorGlobalRulesFile)
	err = os.WriteFile(globalFile, []byte("# Test global rules"), 0644)
	assert.NoError(t, err)

	// Set the TEMPLATE_DIR environment variable
	oldTemplateDir := os.Getenv("TEMPLATE_DIR")
	defer os.Setenv("TEMPLATE_DIR", oldTemplateDir)
	os.Setenv("TEMPLATE_DIR", tempDir)

	// Create a temporary working directory
	workDir := t.TempDir()
	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(oldWd)
	err = os.Chdir(workDir)
	assert.NoError(t, err)

	// Create the installer
	installer, err := NewCursorInstaller(English)
	assert.NoError(t, err)

	// Call the function being tested
	err = installer.Install(Global)
	assert.NoError(t, err)

	// Check that the file was installed
	destPath := filepath.Join(workDir, CursorGlobalRulesFile)
	assert.FileExists(t, destPath)

	// Clean up the installed file
	defer os.Remove(destPath)
	defer os.Remove(destPath + ".backup")

	// Check the content of the installed file
	content, err := os.ReadFile(destPath)
	assert.NoError(t, err)
	assert.Equal(t, "# Test global rules", string(content))
}

// TestCursorInstallerInstallAll tests the Install function with All installation type
func TestCursorInstallerInstallAll(t *testing.T) {
	// Set up a temporary template directory
	tempDir := t.TempDir()
	cursorDir := filepath.Join(tempDir, "cursor")
	err := os.MkdirAll(filepath.Join(cursorDir, "local"), 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Join(cursorDir, "global"), 0755)
	assert.NoError(t, err)

	// Create test template files
	localFile := filepath.Join(cursorDir, "local", CursorLocalRulesFile)
	err = os.WriteFile(localFile, []byte("# Test local rules"), 0644)
	assert.NoError(t, err)

	globalFile := filepath.Join(cursorDir, "global", CursorGlobalRulesFile)
	err = os.WriteFile(globalFile, []byte("# Test global rules"), 0644)
	assert.NoError(t, err)

	// Set the TEMPLATE_DIR environment variable
	oldTemplateDir := os.Getenv("TEMPLATE_DIR")
	defer os.Setenv("TEMPLATE_DIR", oldTemplateDir)
	os.Setenv("TEMPLATE_DIR", tempDir)

	// Create a temporary working directory
	workDir := t.TempDir()
	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(oldWd)
	err = os.Chdir(workDir)
	assert.NoError(t, err)

	// Create the installer
	installer, err := NewCursorInstaller(English)
	assert.NoError(t, err)

	// Call the function being tested
	err = installer.Install(All)
	assert.NoError(t, err)

	// Check that the local file was installed
	localDestPath := filepath.Join(workDir, CursorLocalRulesDir, CursorLocalRulesFile)
	assert.FileExists(t, localDestPath)

	// Check the content of the installed local file
	localContent, err := os.ReadFile(localDestPath)
	assert.NoError(t, err)
	assert.Equal(t, "# Test local rules", string(localContent))

	// Check that the global file was installed
	globalDestPath := filepath.Join(workDir, CursorGlobalRulesFile)
	assert.FileExists(t, globalDestPath)

	// Clean up the installed global file
	defer os.Remove(globalDestPath)
	defer os.Remove(globalDestPath + ".backup")

	// Check the content of the installed global file
	globalContent, err := os.ReadFile(globalDestPath)
	assert.NoError(t, err)
	assert.Equal(t, "# Test global rules", string(globalContent))
}
