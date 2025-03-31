package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

// InstallType represents the type of installation.
type InstallType int

const (
	// Local represents local installation.
	Local InstallType = iota
	// Global represents global installation.
	Global
	// All represents both local and global installation.
	All
)

// String returns the string representation of InstallType.
func (t InstallType) String() string {
	switch t {
	case Local:
		return "local"
	case Global:
		return "global"
	case All:
		return "all"
	default:
		return "unknown"
	}
}

// EditorConfig represents the configuration for an editor.
type EditorConfig struct {
	Name            string
	GlobalSupported bool
	LocalPath       string
	GlobalPath      string
	LocalFileName   string
	GlobalFileName  string
}

// GetRuleFilePaths returns the rule file paths for the specified mode.
func (c *EditorConfig) GetRuleFilePaths(mode string) ([]string, error) {
	switch mode {
	case "local":
		return []string{filepath.Join(c.LocalPath, c.LocalFileName)}, nil
	case "global":
		if !c.GlobalSupported {
			return nil, fmt.Errorf("global mode not supported for editor %s", c.Name)
		}

		return []string{filepath.Join(c.GlobalPath, c.GlobalFileName)}, nil
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}
}

// editorConfigs maps editor names to their configurations.
var editorConfigs = map[string]func() (EditorConfig, error){
	"windsurf": func() (EditorConfig, error) {
		home, err := homedir.Dir()
		if err != nil {
			return EditorConfig{}, fmt.Errorf("failed to get home directory: %w", err)
		}

		var globalDestDir string
		switch runtime.GOOS {
		case "darwin", "linux", "windows":
			globalDestDir = filepath.Join(home, ".codeium", "windsurf", "memories")
		default:
			return EditorConfig{}, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
		}

		return EditorConfig{
			LocalPath:       ".",
			GlobalPath:      globalDestDir,
			GlobalSupported: true,
		}, nil
	},
	"cursor": func() (EditorConfig, error) {
		// Store local rules in the ./.cursor/rules/ directory
		// Also store global rules in the same directory (based on user request)
		localDestDir := filepath.Join(".", ".cursor", "rules")

		return EditorConfig{
			LocalPath:       localDestDir,
			GlobalPath:      localDestDir, // Use the same directory as local rules
			GlobalSupported: true,         // Support global rules
		}, nil
	},
}

// FileSystem interface defines file system operations.
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	CopyFile(src, dest string) error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	Stat(path string) (os.FileInfo, error)
	Rename(oldpath, newpath string) error
}

// DefaultFileSystem implements FileSystem interface using OS operations.
type DefaultFileSystem struct{}

func (fs *DefaultFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *DefaultFileSystem) CopyFile(src, dest string) error {
	return CopyFile(src, dest)
}

func (fs *DefaultFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs *DefaultFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (fs *DefaultFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (fs *DefaultFileSystem) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// CopyFile copies a file from src to dest.
func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer func(srcFile *os.File) {
		if closeErr := srcFile.Close(); closeErr != nil {
			fmt.Printf("warning: could not close source file: %v\n", closeErr)
		}
	}(srcFile)

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer func(destFile *os.File) {
		if closeErr := destFile.Close(); closeErr != nil {
			fmt.Printf("warning: could not close destination file: %v\n", closeErr)
		}
	}(destFile)

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("could not get source file info: %w", err)
	}

	err = os.Chmod(dest, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to set permissions on destination file: %w", err)
	}

	return nil
}

// GetSupportedEditors returns a list of supported editors.
func GetSupportedEditors() []string {
	editors := make([]string, 0, len(editorConfigs))
	for editor := range editorConfigs {
		editors = append(editors, editor)
	}

	return editors
}

// IsEditorSupported checks if an editor is supported.
func IsEditorSupported(editor string) bool {
	_, ok := editorConfigs[editor]

	return ok
}

// IsGlobalModeSupported checks if the global mode is supported for the editor.
func IsGlobalModeSupported(editor string) bool {
	configFn, ok := editorConfigs[editor]
	if !ok {
		return false
	}

	config, err := configFn()
	if err != nil {
		return false
	}

	return config.GlobalSupported
}

// Install installs rules for the specified editor and installation type.
func Install(editor string, installType InstallType) error {
	return InstallWithKey(editor, installType, "default")
}

// validateInstallParams validates installation parameters.
func validateInstallParams(editor string, installType InstallType, key string) error {
	if !IsEditorSupported(editor) {
		return fmt.Errorf("unsupported editor: %s", editor)
	}

	if key == "" {
		return fmt.Errorf("key is required")
	}

	if installType == Global && !IsGlobalModeSupported(editor) {
		return fmt.Errorf("editor '%s' does not support global mode installation", editor)
	}

	return nil
}

// getRulePaths gets rule paths based on installation type.
func getRulePaths(fs FileSystem, editorConfig *EditorConfig, installType InstallType) ([]string, error) {
	var localRulePaths, globalRulePaths []string
	var err error

	if installType == Local || installType == All {
		localRulePaths, err = getLocalRulePaths(fs, editorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get local rule paths: %w", err)
		}
	}

	if (installType == Global || installType == All) && IsGlobalModeSupported(editorConfig.Name) {
		globalRulePaths, err = getGlobalRulePaths(fs, editorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get global rule paths: %w", err)
		}
	}

	rulePaths := make([]string, len(localRulePaths)+len(globalRulePaths))
	copy(rulePaths, localRulePaths)
	copy(rulePaths[len(localRulePaths):], globalRulePaths)

	return rulePaths, nil
}

// InstallWithKey installs rules for the specified editor with a given key.
func InstallWithKey(editor string, installType InstallType, key string) error {
	if err := validateInstallParams(editor, installType, key); err != nil {
		return err
	}

	fs := NewOsFS()
	editorConfig, err := GetEditorConfig(editor)
	if err != nil {
		return fmt.Errorf("failed to get editor config: %w", err)
	}

	rulePaths, err := getRulePaths(fs, &editorConfig, installType)
	if err != nil {
		return err
	}

	if len(rulePaths) == 0 {
		return fmt.Errorf("no rules found for editor '%s'", editor)
	}

	err = installLocal(fs, &editorConfig, rulePaths)
	if err != nil {
		return fmt.Errorf("failed to install rules: %w", err)
	}

	return nil
}

// createBackup makes a backup of an existing file.
func createBackup(fs FileSystem, filePath string) error {
	// Check if the file exists
	_, err := fs.Stat(filePath)
	if os.IsNotExist(err) {
		// No backup needed if the file doesn't exist
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	// Generate backup filename (original filename + .backup_YYYYMMDDhhmmss)
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup_%s", filePath, timestamp)

	// Rename the file
	err = fs.Rename(filePath, backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	fmt.Printf("Created backup of existing file at: %s\n", backupPath)

	return nil
}

// installLocal installs local rules using the specified rule files.
func installLocal(fs FileSystem, config *EditorConfig, rulePaths []string) error {
	destPath := filepath.Join(config.LocalPath, config.LocalFileName)
	destDir := filepath.Dir(destPath)

	if err := fs.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Create a backup of the existing file if it exists
	if err := createBackup(fs, destPath); err != nil {
		return err
	}

	return combineAndWriteRules(fs, rulePaths, destPath)
}

// combineAndWriteRules combines multiple rule files and writes them to the destination.
func combineAndWriteRules(fs FileSystem, rulePaths []string, destPath string) error {
	var combinedContent strings.Builder

	for _, path := range rulePaths {
		content, err := fs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read rule file '%s': %w", path, err)
		}

		// Add file content with a separator
		if combinedContent.Len() > 0 {
			combinedContent.WriteString("\n\n")
		}
		combinedContent.WriteString(fmt.Sprintf("// From %s\n", filepath.Base(path)))
		combinedContent.Write(content)
	}

	// Write combined content to destination
	if err := fs.WriteFile(destPath, []byte(combinedContent.String()), 0o644); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", destPath, err)
	}

	return nil
}

// getLocalRulePaths returns the paths of local rules.
func getLocalRulePaths(fs FileSystem, config *EditorConfig) ([]string, error) {
	paths, err := config.GetRuleFilePaths("local")
	if err != nil {
		return nil, fmt.Errorf("failed to get local rule paths: %w", err)
	}

	return paths, nil
}

// getGlobalRulePaths returns the paths of global rules.
func getGlobalRulePaths(fs FileSystem, config *EditorConfig) ([]string, error) {
	if !config.GlobalSupported {
		return nil, nil
	}
	paths, err := config.GetRuleFilePaths("global")
	if err != nil {
		return nil, fmt.Errorf("failed to get global rule paths: %w", err)
	}

	return paths, nil
}

// NewOsFS creates a new OS file system implementation.
func NewOsFS() FileSystem {
	return &DefaultFileSystem{}
}

// GetEditorConfig returns the configuration for the specified editor.
func GetEditorConfig(editor string) (EditorConfig, error) {
	configFn, ok := editorConfigs[editor]
	if !ok {
		return EditorConfig{}, fmt.Errorf("unsupported editor: %s", editor)
	}

	return configFn()
}
