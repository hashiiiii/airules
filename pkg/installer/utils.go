package installer

import (
	"fmt"
	"io"
	"os"
)

// copyFile is a utility function to copy a file
func copyFile(src, dest string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy file
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Get source file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("could not get source file info: %w", err)
	}

	// Set permissions on destination file
	err = os.Chmod(dest, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to set permissions on destination file: %w", err)
	}

	return nil
}
