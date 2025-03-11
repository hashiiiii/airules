package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	CopyFile(src, dest string) error
}

type DefaultFileSystem struct{}

func (fs *DefaultFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *DefaultFileSystem) CopyFile(src, dest string) error {
	return CopyFile(src, dest)
}

type WindsurfInstaller struct {
	templateDir    string
	localDestDir   string
	globalDestDir  string
	localFileName  string
	globalFileName string
	lang           Language
	fs             FileSystem
}

func NewWindsurfInstaller(lang Language) (*WindsurfInstaller, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		return nil, fmt.Errorf("failed to get template directory: TEMPLATE_DIR environment variable is not set")
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
		templateDir:    templateDir,
		localDestDir:   localDestDir,
		globalDestDir:  globalDestDir,
		localFileName:  ".windsurfrules",
		globalFileName: "global_rules.md",
		lang:           lang,
		fs:             &DefaultFileSystem{},
	}, nil
}

func (i *WindsurfInstaller) Install(installType InstallType) error {
	switch installType {
	case Local, Global:
		return i.installCore(installType)
	case All:
		if err := i.installCore(Local); err != nil {
			return fmt.Errorf("failed to install local configuration file: %w", err)
		}

		if err := i.installCore(Global); err != nil {
			return fmt.Errorf("failed to install global configuration file: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unknown installation type: %v", installType)
	}
}

func (i *WindsurfInstaller) installCore(installType InstallType) error {
	srcPath, destPath, destDir, err := GetInstallPath(installType, i.templateDir, i.localDestDir, i.globalDestDir, i.localFileName, i.globalFileName, i.lang)
	if err != nil {
		return err
	}

	if err := i.fs.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Check if the destination file already exists
	if _, err := os.Stat(destPath); err == nil {
		// Create a backup of the existing file
		backupPath := destPath + ".backup"
		if err := i.fs.CopyFile(destPath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup of existing file: %w", err)
		}
		fmt.Printf("Created backup of existing file at %s\n", backupPath)
	}

	if err := i.fs.CopyFile(srcPath, destPath); err != nil {
		return fmt.Errorf("failed to copy template file: %w", err)
	}

	fmt.Printf("Installed %s Windsurf rules to %s\n",
		installType.String(), destPath)

	return nil
}
