package installer

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashiiiii/airules/pkg/remote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFetcher is a mock implementation of remote.Fetcher
type MockFetcher struct {
	mock.Mock
}

// ListRuleSets implements the remote.Fetcher interface
func (m *MockFetcher) ListRuleSets(ctx context.Context) ([]remote.RuleSet, error) {
	args := m.Called(ctx)
	return args.Get(0).([]remote.RuleSet), args.Error(1)
}

// FetchRuleSet implements the remote.Fetcher interface
func (m *MockFetcher) FetchRuleSet(ctx context.Context, ruleSet remote.RuleSet) (io.ReadCloser, error) {
	args := m.Called(ctx, ruleSet)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// FetchRuleSetByName implements the remote.Fetcher interface
func (m *MockFetcher) FetchRuleSetByName(ctx context.Context, name string) (io.ReadCloser, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// TestNewRemoteInstaller tests the NewRemoteInstaller function
func TestNewRemoteInstaller(t *testing.T) {
	// Create a mock fetcher
	mockFetcher := new(MockFetcher)

	// Call the function being tested
	installer := NewRemoteInstaller(mockFetcher)

	// Assert that the installer is not nil
	assert.NotNil(t, installer)
	assert.Equal(t, mockFetcher, installer.fetcher)
}

// TestRemoteInstallerInstallRuleSet tests the InstallRuleSet function
func TestRemoteInstallerInstallRuleSet(t *testing.T) {
	// Create a mock fetcher
	mockFetcher := new(MockFetcher)

	// Create a sample rule set
	ruleSet := remote.RuleSet{
		Name: "cursor-rules-example",
		Path: "rules/cursor-rules-example",
		Type: "cursor",
		URL:  "https://github.com/owner/repo/tree/main/rules/cursor-rules-example",
	}

	// Create a sample rule content
	ruleContent := "# Example Cursor Rules"
	mockResponse := io.NopCloser(strings.NewReader(ruleContent))

	// Set up the mock expectation
	mockFetcher.On("FetchRuleSet", mock.Anything, ruleSet).Return(mockResponse, nil)

	// Create the installer
	installer := NewRemoteInstaller(mockFetcher)

	// Create a temporary working directory
	workDir := t.TempDir()
	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(oldWd)
	err = os.Chdir(workDir)
	assert.NoError(t, err)

	// Call the function being tested
	ctx := context.Background()
	err = installer.InstallRuleSet(ctx, ruleSet, Local)
	assert.NoError(t, err)

	// Check that the file was installed
	destPath := filepath.Join(workDir, CursorLocalRulesDir, CursorLocalRulesFile)
	assert.FileExists(t, destPath)

	// Check the content of the installed file
	content, err := os.ReadFile(destPath)
	assert.NoError(t, err)
	assert.Equal(t, ruleContent, string(content))

	// Assert that the mock expectations were met
	mockFetcher.AssertExpectations(t)
}

// TestRemoteInstallerInstallRuleSetByName tests the InstallRuleSetByName function
func TestRemoteInstallerInstallRuleSetByName(t *testing.T) {
	// Create a mock fetcher
	mockFetcher := new(MockFetcher)

	// Create a sample rule set
	ruleSets := []remote.RuleSet{
		{
			Name: "cursor-rules-example",
			Path: "rules/cursor-rules-example",
			Type: "cursor",
			URL:  "https://github.com/owner/repo/tree/main/rules/cursor-rules-example",
		},
	}

	// Create a sample rule content
	ruleContent := "# Example Cursor Rules"
	mockResponse := io.NopCloser(strings.NewReader(ruleContent))

	// Set up the mock expectations
	mockFetcher.On("FetchRuleSetByName", mock.Anything, "cursor-rules-example").Return(mockResponse, nil)
	mockFetcher.On("ListRuleSets", mock.Anything).Return(ruleSets, nil)

	// Add mock for FetchRuleSet which is called by InstallRuleSet
	mockFetcher.On("FetchRuleSet", mock.Anything, ruleSets[0]).Return(mockResponse, nil)

	// Create the installer
	installer := NewRemoteInstaller(mockFetcher)

	// Create a temporary working directory
	workDir := t.TempDir()
	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(oldWd)
	err = os.Chdir(workDir)
	assert.NoError(t, err)

	// Call the function being tested
	ctx := context.Background()
	err = installer.InstallRuleSetByName(ctx, "cursor-rules-example", Local)
	assert.NoError(t, err)

	// Check that the file was installed
	destPath := filepath.Join(workDir, ".cursor", "rules", "project_rules.mdc")
	assert.FileExists(t, destPath)

	// Check the content of the installed file
	content, err := os.ReadFile(destPath)
	assert.NoError(t, err)
	assert.Equal(t, ruleContent, string(content))

	// Assert that the mock expectations were met
	mockFetcher.AssertExpectations(t)
}

// TestRemoteInstallerListRuleSets tests the ListRuleSets function
func TestRemoteInstallerListRuleSets(t *testing.T) {
	// Create a mock fetcher
	mockFetcher := new(MockFetcher)

	// Create a sample rule sets
	ruleSets := []remote.RuleSet{
		{
			Name: "cursor-rules-example",
			Path: "rules/cursor-rules-example",
			Type: "cursor",
			URL:  "https://github.com/owner/repo/tree/main/rules/cursor-rules-example",
		},
		{
			Name: "windsurf-rules-example",
			Path: "rules/windsurf-rules-example",
			Type: "windsurf",
			URL:  "https://github.com/owner/repo/tree/main/rules/windsurf-rules-example",
		},
	}

	// Set up the mock expectation
	mockFetcher.On("ListRuleSets", mock.Anything).Return(ruleSets, nil)

	// Create the installer
	installer := NewRemoteInstaller(mockFetcher)

	// Call the function being tested
	ctx := context.Background()
	result, err := installer.ListRuleSets(ctx)

	// Assert that there was no error
	assert.NoError(t, err)
	assert.Equal(t, ruleSets, result)

	// Assert that the mock expectations were met
	mockFetcher.AssertExpectations(t)
}
