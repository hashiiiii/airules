package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashiiiii/airules/pkg/config"
	"github.com/mitchellh/go-homedir"
)

type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	CopyFile(src, dest string) error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
}

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

type WindsurfInstaller struct {
	localDestDir   string
	globalDestDir  string
	localFileName  string
	globalFileName string
	fs             FileSystem
}

func NewWindsurfInstaller() (*WindsurfInstaller, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	var localDestDir, globalDestDir string

	switch runtime.GOOS {
	case "darwin", "linux", "windows":
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(home, ".codeium", "windsurf", "memories")
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return &WindsurfInstaller{
		localDestDir:   localDestDir,
		globalDestDir:  globalDestDir,
		localFileName:  ".windsurfrules",
		globalFileName: "global_rules.md",
		fs:             &DefaultFileSystem{},
	}, nil
}

// Install installs rules based on the specified key
func (i *WindsurfInstaller) Install(installType InstallType) error {
	return i.InstallWithKey(installType, "default")
}

// InstallWithKey installs rules based on the specified key
func (i *WindsurfInstaller) InstallWithKey(installType InstallType, key string) error {
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
		rulePaths, err = config.GetRuleFilePaths("windsurf", "local", key)
	case Global:
		rulePaths, err = config.GetRuleFilePaths("windsurf", "global", key)
	case All:
		// For "all", we'll combine both local and global rules
		localRulePaths, localErr := config.GetRuleFilePaths("windsurf", "local", key)
		if localErr != nil {
			return fmt.Errorf("failed to get local rule file paths: %w", localErr)
		}

		globalRulePaths, globalErr := config.GetRuleFilePaths("windsurf", "global", key)
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
func (i *WindsurfInstaller) installLocal(rulePaths []string) error {
	destPath := filepath.Join(i.localDestDir, i.localFileName)
	destDir := filepath.Dir(destPath)

	if err := i.fs.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	return i.combineAndWriteRules(rulePaths, destPath)
}

// installGlobal installs global rules using the specified rule files
func (i *WindsurfInstaller) installGlobal(rulePaths []string) error {
	destPath := filepath.Join(i.globalDestDir, i.globalFileName)
	destDir := filepath.Dir(destPath)

	if err := i.fs.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	return i.combineAndWriteRules(rulePaths, destPath)
}

// combineAndWriteRules combines multiple rule files and writes them to the destination
func (i *WindsurfInstaller) combineAndWriteRules(rulePaths []string, destPath string) error {
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
