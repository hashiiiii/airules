package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// Language represents the language of the template file
type Language int

const (
	// English represents English language
	English Language = iota
	// Japanese represents Japanese language
	Japanese
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

// String returns the string representation of Language
func (l Language) String() string {
	switch l {
	case English:
		return "en"
	case Japanese:
		return "ja"
	default:
		return "unknown"
	}
}

// Installer interface defines methods that must be implemented by installers for each editor
type Installer interface {
	// Install installs configuration files based on the installation type
	Install(installType InstallType) error
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

func GetInstallPath(installType InstallType, templateDir, localDestDir, globalDestDir, localFileName, globalFileName string, lang Language) (srcPath, destPath, destDir string, err error) {
	switch installType {
	case Local:
		destDir = localDestDir
		fileName := localFileName

		templateFileName := fileName
		if lang == Japanese {
			templateFileName = fileName + "_JA"
		}

		srcPath = filepath.Join(templateDir, "local", templateFileName)
		destPath = filepath.Join(destDir, fileName)
	case Global:
		destDir = globalDestDir
		fileName := globalFileName

		templateFileName := fileName
		if lang == Japanese {
			templateFileName = fileName + "_JA"
		}

		srcPath = filepath.Join(templateDir, "global", templateFileName)
		destPath = filepath.Join(destDir, fileName)
	default:
		return "", "", "", fmt.Errorf("unknown installation type: %v", installType)
	}

	return srcPath, destPath, destDir, nil
}
