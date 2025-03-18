package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashiiiii/airules/pkg/config"
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

// Installer interface defines methods that must be implemented by installers for each editor
type Installer interface {
	// Install installs configuration files based on the installation type
	Install(installType InstallType) error
	// InstallWithKey installs configuration files based on the installation type and key
	InstallWithKey(installType InstallType, key string) error
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

func GetInstallPath(installType InstallType, templateDir, localDestDir, globalDestDir, localFileName, globalFileName string) (srcPath, destPath, destDir string, err error) {
	switch installType {
	case Local:
		destDir = localDestDir
		fileName := localFileName
		srcPath = filepath.Join(templateDir, "local", fileName)
		destPath = filepath.Join(destDir, fileName)
	case Global:
		destDir = globalDestDir
		fileName := globalFileName
		srcPath = filepath.Join(templateDir, "global", fileName)
		destPath = filepath.Join(destDir, fileName)
	default:
		return "", "", "", fmt.Errorf("unknown installation type: %v", installType)
	}

	return srcPath, destPath, destDir, nil
}

// GetInstaller returns an installer for the specified editor
func GetInstaller(editor string) (Installer, error) {
	switch editor {
	case "windsurf":
		return NewWindsurfInstaller()
	case "cursor":
		return NewCursorInstaller()
	default:
		return nil, fmt.Errorf("unsupported editor: %s", editor)
	}
}

// GetSupportedEditors returns a list of supported editors
func GetSupportedEditors() []string {
	return config.GetSupportedEditors()
}

// FileSystem interface defines file system operations
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	CopyFile(src, dest string) error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
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

// BaseInstaller provides common functionality for all installers
type BaseInstaller struct {
	localDestDir   string
	globalDestDir  string
	localFileName  string
	globalFileName string
	editor         string
	fs             FileSystem
}

// Install installs rules based on the specified key
func (i *BaseInstaller) Install(installType InstallType) error {
	return i.InstallWithKey(installType, "default")
}

// InstallWithKey installs rules based on the specified key
func (i *BaseInstaller) InstallWithKey(installType InstallType, key string) error {
	// Ensure config directory exists
	if _, err := config.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to ensure config directory: %w", err)
	}

	// Get rule file paths for the key
	var rulePaths []string
	var err error

	// Get rule file paths based on installation type
	switch installType {
	case Local:
		rulePaths, err = config.GetRuleFilePaths(i.editor, "local", key)
	case Global:
		rulePaths, err = config.GetRuleFilePaths(i.editor, "global", key)
	case All:
		// For "all", we'll combine both local and global rules
		localRulePaths, localErr := config.GetRuleFilePaths(i.editor, "local", key)
		if localErr != nil {
			return fmt.Errorf("failed to get local rule file paths: %w", localErr)
		}

		globalRulePaths, globalErr := config.GetRuleFilePaths(i.editor, "global", key)
		if globalErr != nil {
			return fmt.Errorf("failed to get global rule file paths: %w", globalErr)
		}

		rulePaths = append(localRulePaths, globalRulePaths...)
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

	// Install based on type
	switch installType {
	case Local:
		err = i.installLocal(rulePaths)
	case Global:
		err = i.installGlobal(rulePaths)
	case All:
		if err = i.installLocal(rulePaths); err != nil {
			return err
		}
		err = i.installGlobal(rulePaths)
	default:
		return fmt.Errorf("unknown install type: %v", installType)
	}

	return err
}

// installLocal installs local rules using the specified rule files
func (i *BaseInstaller) installLocal(rulePaths []string) error {
	destPath := filepath.Join(i.localDestDir, i.localFileName)
	destDir := filepath.Dir(destPath)

	if err := i.fs.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	return i.combineAndWriteRules(rulePaths, destPath)
}

// installGlobal installs global rules using the specified rule files
func (i *BaseInstaller) installGlobal(rulePaths []string) error {
	destPath := filepath.Join(i.globalDestDir, i.globalFileName)
	destDir := filepath.Dir(destPath)

	if err := i.fs.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	return i.combineAndWriteRules(rulePaths, destPath)
}

// combineAndWriteRules combines multiple rule files and writes them to the destination
func (i *BaseInstaller) combineAndWriteRules(rulePaths []string, destPath string) error {
	var combinedContent strings.Builder

	for _, path := range rulePaths {
		content, err := i.fs.ReadFile(path)
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
	if err := i.fs.WriteFile(destPath, []byte(combinedContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", destPath, err)
	}

	return nil
}
