package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hashiiiii/airules/pkg/config"
	"github.com/mitchellh/go-homedir"
)

// InstallType represents the type of installation
type InstallType int

const (
	// Local represents local installation
	Local InstallType = iota
	// Global represents global installation
	Global
	// All represents both local and global installation
	All
)

// String returns the string representation of InstallType
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

// EditorConfig holds the configuration for a specific editor
type EditorConfig struct {
	LocalDestDir    string
	GlobalDestDir   string
	LocalFileName   string
	GlobalFileName  string
	GlobalSupported bool
}

// editorConfigs maps editor names to their configurations
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
			LocalDestDir:    ".",
			GlobalDestDir:   globalDestDir,
			LocalFileName:   ".windsurfrules",
			GlobalFileName:  "global_rules.md",
			GlobalSupported: true,
		}, nil
	},
	"cursor": func() (EditorConfig, error) {
		// Store local rules in the ./.cursor/rules/ directory
		// Also store global rules in the same directory (based on user request)
		localDestDir := filepath.Join(".", ".cursor", "rules")
		return EditorConfig{
			LocalDestDir:    localDestDir,
			GlobalDestDir:   localDestDir, // Use the same directory as local rules
			LocalFileName:   "project_rules.mdc",
			GlobalFileName:  "global_rules.mdc",
			GlobalSupported: true, // Support global rules
		}, nil
	},
}

// FileSystem interface defines file system operations
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	CopyFile(src, dest string) error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	Stat(path string) (os.FileInfo, error)
	Rename(oldpath, newpath string) error
}

// DefaultFileSystem implements FileSystem interface using OS operations
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

// CopyFile copies a file from src to dest
func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			_ = fmt.Errorf("could not close source file: %w", err)
		}
	}(srcFile)

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			_ = fmt.Errorf("could not close source file: %w", err)
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

// GetSupportedEditors returns a list of supported editors
func GetSupportedEditors() []string {
	editors := make([]string, 0, len(editorConfigs))
	for editor := range editorConfigs {
		editors = append(editors, editor)
	}
	return editors
}

// IsEditorSupported checks if an editor is supported
func IsEditorSupported(editor string) bool {
	_, ok := editorConfigs[editor]
	return ok
}

// IsGlobalModeSupported checks if the global mode is supported for the editor
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

// Install installs rules for the specified editor and installation type
func Install(editor string, installType InstallType) error {
	return InstallWithKey(editor, installType, "default")
}

// InstallWithKey installs rules for the specified editor, installation type, and template key
func InstallWithKey(editor string, installType InstallType, key string) error {
	// Get editor configuration
	configFn, ok := editorConfigs[editor]
	if !ok {
		return fmt.Errorf("unsupported editor: %s", editor)
	}

	editorConfig, err := configFn()
	if err != nil {
		return fmt.Errorf("failed to get editor configuration: %w", err)
	}

	// If global mode installation is requested but the editor doesn't support it
	if (installType == Global || installType == All) && !editorConfig.GlobalSupported {
		if installType == Global {
			return fmt.Errorf("editor '%s' does not support global mode installation through files", editor)
		}
		// For 'All' type, show warning and install only local rules
		fmt.Printf("Warning: Editor '%s' does not support global mode installation through files. Only local rules will be installed.\n", editor)
		installType = Local
	}

	// Ensure config directory exists
	if _, err := config.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to ensure config directory: %w", err)
	}

	// Get rule file paths for the key
	var rulePaths []string

	// Get rule file paths based on installation type
	switch installType {
	case Local:
		rulePaths, err = config.GetRuleFilePaths(editor, "local", key)
	case Global:
		rulePaths, err = config.GetRuleFilePaths(editor, "global", key)
	case All:
		// For "all", we'll combine both local and global rules
		localRulePaths, localErr := config.GetRuleFilePaths(editor, "local", key)
		if localErr != nil {
			return fmt.Errorf("failed to get local rule file paths: %w", localErr)
		}

		if editorConfig.GlobalSupported {
			globalRulePaths, globalErr := config.GetRuleFilePaths(editor, "global", key)
			if globalErr != nil {
				return fmt.Errorf("failed to get global rule file paths: %w", globalErr)
			}
			rulePaths = append(localRulePaths, globalRulePaths...)
		} else {
			rulePaths = localRulePaths
		}
		err = nil
	default:
		return fmt.Errorf("unknown install type: %v", installType)
	}

	if err != nil {
		return fmt.Errorf("failed to get rule file paths: %w", err)
	}

	// Check if rule files exist
	for _, path := range rulePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("rule file '%s' not found", path)
		}
	}

	fs := &DefaultFileSystem{}

	// Install based on type
	switch installType {
	case Local:
		err = installLocal(fs, editorConfig, rulePaths)
	case Global:
		err = installGlobal(fs, editorConfig, rulePaths)
	case All:
		if err = installLocal(fs, editorConfig, rulePaths); err != nil {
			return err
		}
		if editorConfig.GlobalSupported {
			err = installGlobal(fs, editorConfig, rulePaths)
		}
	default:
		return fmt.Errorf("unknown install type: %v", installType)
	}

	return err
}

// createBackup makes a backup of an existing file
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

// installLocal installs local rules using the specified rule files
func installLocal(fs FileSystem, config EditorConfig, rulePaths []string) error {
	destPath := filepath.Join(config.LocalDestDir, config.LocalFileName)
	destDir := filepath.Dir(destPath)

	if err := fs.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Create a backup of the existing file if it exists
	if err := createBackup(fs, destPath); err != nil {
		return err
	}

	return combineAndWriteRules(fs, rulePaths, destPath)
}

// installGlobal installs global rules using the specified rule files
func installGlobal(fs FileSystem, config EditorConfig, rulePaths []string) error {
	if !config.GlobalSupported || config.GlobalDestDir == "" || config.GlobalFileName == "" {
		return fmt.Errorf("global rules installation not supported for this editor")
	}

	destPath := filepath.Join(config.GlobalDestDir, config.GlobalFileName)
	destDir := filepath.Dir(destPath)

	if err := fs.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Create a backup of the existing file if it exists
	if err := createBackup(fs, destPath); err != nil {
		return err
	}

	return combineAndWriteRules(fs, rulePaths, destPath)
}

// combineAndWriteRules combines multiple rule files and writes them to the destination
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
	if err := fs.WriteFile(destPath, []byte(combinedContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", destPath, err)
	}

	return nil
}
