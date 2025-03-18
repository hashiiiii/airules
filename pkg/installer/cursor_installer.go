package installer

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

// CursorInstaller implements Installer interface for Cursor editor
type CursorInstaller struct {
	BaseInstaller
}

func NewCursorInstaller() (*CursorInstaller, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	var localDestDir, globalDestDir string

	switch runtime.GOOS {
	case "darwin", "linux":
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(home, ".cursor", "ai")
	case "windows":
		localDestDir = filepath.Join(".")
		globalDestDir = filepath.Join(home, "AppData", "Roaming", "cursor", "ai")
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return &CursorInstaller{
		BaseInstaller: BaseInstaller{
			localDestDir:   localDestDir,
			globalDestDir:  globalDestDir,
			localFileName:  "project_rules.mdc",
			globalFileName: "global_rules.mdc",
			editor:         "cursor",
			fs:             &DefaultFileSystem{},
		},
	}, nil
}

// Install delegates to BaseInstaller.Install
func (i *CursorInstaller) Install(installType InstallType) error {
	return i.BaseInstaller.Install(installType)
}

// InstallWithKey delegates to BaseInstaller.InstallWithKey
func (i *CursorInstaller) InstallWithKey(installType InstallType, key string) error {
	return i.BaseInstaller.InstallWithKey(installType, key)
}
