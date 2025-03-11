package installer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashiiiii/airules/pkg/remote"
)

// RemoteInstaller installs rule sets from remote repositories
type RemoteInstaller struct {
	fetcher remote.Fetcher
}

// NewRemoteInstaller creates a new RemoteInstaller
func NewRemoteInstaller(fetcher remote.Fetcher) *RemoteInstaller {
	return &RemoteInstaller{
		fetcher: fetcher,
	}
}

// ListRuleSets lists available rule sets from the remote repository
func (i *RemoteInstaller) ListRuleSets(ctx context.Context) ([]remote.RuleSet, error) {
	return i.fetcher.ListRuleSets(ctx)
}

// InstallRuleSet installs a rule set from the remote repository
func (i *RemoteInstaller) InstallRuleSet(ctx context.Context, ruleSet remote.RuleSet, installType InstallType) error {
	// Fetch the rule set content
	content, err := i.fetcher.FetchRuleSet(ctx, ruleSet)
	if err != nil {
		return fmt.Errorf("failed to fetch rule set: %w", err)
	}
	defer content.Close()

	// Determine the destination path based on the rule set type and installation type
	var destPath string
	switch ruleSet.Type {
	case "cursor":
		destPath, err = i.getCursorDestPath(installType)
	case "windsurf":
		destPath, err = i.getWindsurfDestPath(installType)
	default:
		return fmt.Errorf("unsupported rule set type: %s", ruleSet.Type)
	}
	if err != nil {
		return fmt.Errorf("failed to determine destination path: %w", err)
	}

	// Create the destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Check if the destination file already exists
	if _, err := os.Stat(destPath); err == nil {
		// Create a backup of the existing file
		backupPath := destPath + ".bak"
		if err := os.Rename(destPath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup of existing file: %w", err)
		}
		fmt.Printf("Created backup of existing file at %s\n", backupPath)
	}

	// Create the destination file
	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	// Copy the content to the destination file
	if _, err := io.Copy(dest, content); err != nil {
		return fmt.Errorf("failed to write content to destination file: %w", err)
	}

	fmt.Printf("Installed rule set to %s\n", destPath)
	return nil
}

// getCursorDestPath returns the destination path for a Cursor rule set
func (i *RemoteInstaller) getCursorDestPath(installType InstallType) (string, error) {
	switch installType {
	case Local:
		// Create the .cursor/rules directory if it doesn't exist
		cursorDir := filepath.Join(".", ".cursor", "rules")
		if err := os.MkdirAll(cursorDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create .cursor/rules directory: %w", err)
		}
		return filepath.Join(cursorDir, "project_rules.mdc"), nil
	case Global:
		// Get the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
		return filepath.Join(cwd, ".cursorrules"), nil
	default:
		return "", fmt.Errorf("unsupported installation type: %s", installType)
	}
}

// getWindsurfDestPath returns the destination path for a Windsurf rule set
func (i *RemoteInstaller) getWindsurfDestPath(installType InstallType) (string, error) {
	switch installType {
	case Local:
		return filepath.Join(".", ".windsurfrules"), nil
	case Global:
		// Get the user's home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user's home directory: %w", err)
		}
		// Create the ~/.codeium/windsurf/memories directory if it doesn't exist
		windsurfDir := filepath.Join(home, ".codeium", "windsurf", "memories")
		if err := os.MkdirAll(windsurfDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create ~/.codeium/windsurf/memories directory: %w", err)
		}
		return filepath.Join(windsurfDir, "global_rules.md"), nil
	default:
		return "", fmt.Errorf("unsupported installation type: %s", installType)
	}
}

// InstallRuleSetByName installs a rule set by name
func (i *RemoteInstaller) InstallRuleSetByName(ctx context.Context, name string, installType InstallType) error {
	// Fetch the rule set
	ruleContent, err := i.fetcher.FetchRuleSetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to fetch rule set: %w", err)
	}
	defer ruleContent.Close()

	// List all rule sets to find the one with the given name
	ruleSets, err := i.ListRuleSets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list rule sets: %w", err)
	}

	// Find the rule set with the given name
	var targetRuleSet remote.RuleSet
	for _, ruleSet := range ruleSets {
		if ruleSet.Name == name {
			targetRuleSet = ruleSet
			break
		}
	}

	if targetRuleSet.Name == "" {
		return fmt.Errorf("rule set not found: %s", name)
	}

	// Install the rule set
	return i.InstallRuleSet(ctx, targetRuleSet, installType)
}
