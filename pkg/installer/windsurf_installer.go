package installer

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

// WindsurfInstaller implements Installer interface for Windsurf editor
type WindsurfInstaller struct {
	BaseInstaller
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
		BaseInstaller: BaseInstaller{
			localDestDir:   localDestDir,
			globalDestDir:  globalDestDir,
			localFileName:  ".windsurfrules",
			globalFileName: "global_rules.md",
			editor:         "windsurf",
			fs:             &DefaultFileSystem{},
		},
	}, nil
}

// Install delegates to BaseInstaller.Install
func (i *WindsurfInstaller) Install(installType InstallType) error {
	return i.BaseInstaller.Install(installType)
}

// InstallWithKey delegates to BaseInstaller.InstallWithKey
func (i *WindsurfInstaller) InstallWithKey(installType InstallType, key string) error {
	return i.BaseInstaller.InstallWithKey(installType, key)
}
