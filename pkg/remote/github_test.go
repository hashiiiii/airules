package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewGitHubFetcher tests the NewGitHubFetcher function
func TestNewGitHubFetcher(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)

	// Call the function being tested
	fetcher := NewGitHubFetcher(mockClient)

	// Assert that the fetcher is not nil
	assert.NotNil(t, fetcher)
	assert.Equal(t, mockClient, fetcher.client)
	assert.Equal(t, DefaultGitHubAPIURL, fetcher.apiURL)
	assert.Equal(t, DefaultGitHubRawURL, fetcher.rawURL)
	assert.Equal(t, DefaultAwesomeCursorRulesOwner, fetcher.owner)
	assert.Equal(t, DefaultAwesomeCursorRulesRepo, fetcher.repo)
	assert.Equal(t, DefaultAwesomeCursorRulesBranch, fetcher.branch)
	assert.Equal(t, DefaultAwesomeCursorRulesPath, fetcher.rulesPath)
}

// TestNewGitHubFetcherWithOptions tests the NewGitHubFetcher function with options
func TestNewGitHubFetcherWithOptions(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)

	// Define custom values
	customAPIURL := "https://custom-api.github.com"
	customRawURL := "https://custom-raw.githubusercontent.com"
	customOwner := "custom-owner"
	customRepo := "custom-repo"
	customBranch := "custom-branch"
	customRulesPath := "custom-rules-path"

	// Call the function being tested with options
	fetcher := NewGitHubFetcher(
		mockClient,
		WithGitHubAPIURL(customAPIURL),
		WithGitHubRawURL(customRawURL),
		WithGitHubOwner(customOwner),
		WithGitHubRepo(customRepo),
		WithGitHubBranch(customBranch),
		WithGitHubRulesPath(customRulesPath),
	)

	// Assert that the fetcher is not nil and has the custom values
	assert.NotNil(t, fetcher)
	assert.Equal(t, mockClient, fetcher.client)
	assert.Equal(t, customAPIURL, fetcher.apiURL)
	assert.Equal(t, customRawURL, fetcher.rawURL)
	assert.Equal(t, customOwner, fetcher.owner)
	assert.Equal(t, customRepo, fetcher.repo)
	assert.Equal(t, customBranch, fetcher.branch)
	assert.Equal(t, customRulesPath, fetcher.rulesPath)
}

// TestListRuleSets tests the ListRuleSets function
func TestListRuleSets(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)

	// Create a sample response
	contents := []GitHubContent{
		{
			Name:    "cursor-rules-example",
			Path:    "rules/cursor-rules-example",
			Type:    "dir",
			URL:     "https://api.github.com/repos/owner/repo/contents/rules/cursor-rules-example",
			HTMLURL: "https://github.com/owner/repo/tree/main/rules/cursor-rules-example",
		},
		{
			Name:    "windsurf-rules-example",
			Path:    "rules/windsurf-rules-example",
			Type:    "dir",
			URL:     "https://api.github.com/repos/owner/repo/contents/rules/windsurf-rules-example",
			HTMLURL: "https://github.com/owner/repo/tree/main/rules/windsurf-rules-example",
		},
		{
			Name:    "README.md",
			Path:    "rules/README.md",
			Type:    "file",
			URL:     "https://api.github.com/repos/owner/repo/contents/rules/README.md",
			HTMLURL: "https://github.com/owner/repo/blob/main/rules/README.md",
		},
	}

	// Marshal the contents to JSON
	contentsJSON, err := json.Marshal(contents)
	assert.NoError(t, err)

	// Set up the mock response
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(contentsJSON)),
	}

	// Set up the mock expectation
	mockClient.On("Do", mock.Anything).Return(response, nil)

	// Create a fetcher
	fetcher := NewGitHubFetcher(mockClient)

	// Call the function being tested
	ctx := context.Background()
	ruleSets, err := fetcher.ListRuleSets(ctx)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the rule sets were parsed correctly
	assert.Len(t, ruleSets, 2)
	assert.Equal(t, "cursor-rules-example", ruleSets[0].Name)
	assert.Equal(t, "cursor", ruleSets[0].Type)
	assert.Equal(t, "windsurf-rules-example", ruleSets[1].Name)
	assert.Equal(t, "windsurf", ruleSets[1].Type)

	// Assert that the mock expectations were met
	mockClient.AssertExpectations(t)
}

// TestFetchRuleSet tests the FetchRuleSet function
func TestFetchRuleSet(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)

	// Create a sample rule set
	ruleSet := RuleSet{
		Name: "cursor-rules-example",
		Path: "rules/cursor-rules-example",
		Type: "cursor",
		URL:  "https://github.com/owner/repo/tree/main/rules/cursor-rules-example",
	}

	// Create a sample directory contents response
	dirContents := []GitHubContent{
		{
			Name:        "example.mdc",
			Path:        "rules/cursor-rules-example/example.mdc",
			Type:        "file",
			DownloadURL: "https://raw.githubusercontent.com/owner/repo/main/rules/cursor-rules-example/example.mdc",
		},
		{
			Name:        "README.md",
			Path:        "rules/cursor-rules-example/README.md",
			Type:        "file",
			DownloadURL: "https://raw.githubusercontent.com/owner/repo/main/rules/cursor-rules-example/README.md",
		},
	}

	// Marshal the directory contents to JSON
	dirContentsJSON, err := json.Marshal(dirContents)
	assert.NoError(t, err)

	// Set up the first mock response (directory contents)
	dirResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(dirContentsJSON)),
	}

	// Set up the second mock response (file contents)
	fileContents := "# Example Cursor Rules"
	fileResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(fileContents)),
	}

	// Set up the mock expectations
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/repos/PatrickJS/awesome-cursorrules/contents/rules/cursor-rules-example"
	})).Return(dirResponse, nil)
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == "https://raw.githubusercontent.com/owner/repo/main/rules/cursor-rules-example/example.mdc"
	})).Return(fileResponse, nil)

	// Create a fetcher
	fetcher := NewGitHubFetcher(mockClient)

	// Call the function being tested
	ctx := context.Background()
	body, err := fetcher.FetchRuleSet(ctx, ruleSet)

	// Assert that there was no error
	assert.NoError(t, err)

	// Read the response body
	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, string(data))

	// Assert that the mock expectations were met
	mockClient.AssertExpectations(t)
}

// TestFetchRuleSetByName tests the FetchRuleSetByName function
func TestFetchRuleSetByName(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)

	// Create a sample list response
	listContents := []GitHubContent{
		{
			Name:    "cursor-rules-example",
			Path:    "rules/cursor-rules-example",
			Type:    "dir",
			URL:     "https://api.github.com/repos/owner/repo/contents/rules/cursor-rules-example",
			HTMLURL: "https://github.com/owner/repo/tree/main/rules/cursor-rules-example",
		},
	}

	// Marshal the list contents to JSON
	listContentsJSON, err := json.Marshal(listContents)
	assert.NoError(t, err)

	// Set up the first mock response (list contents)
	listResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(listContentsJSON)),
	}

	// Create a sample directory contents response
	dirContents := []GitHubContent{
		{
			Name:        "example.mdc",
			Path:        "rules/cursor-rules-example/example.mdc",
			Type:        "file",
			DownloadURL: "https://raw.githubusercontent.com/owner/repo/main/rules/cursor-rules-example/example.mdc",
		},
	}

	// Marshal the directory contents to JSON
	dirContentsJSON, err := json.Marshal(dirContents)
	assert.NoError(t, err)

	// Set up the second mock response (directory contents)
	dirResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(dirContentsJSON)),
	}

	// Set up the third mock response (file contents)
	fileContents := "# Example Cursor Rules"
	fileResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(fileContents)),
	}

	// Set up the mock expectations
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/repos/PatrickJS/awesome-cursorrules/contents/rules"
	})).Return(listResponse, nil)
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/repos/PatrickJS/awesome-cursorrules/contents/rules/cursor-rules-example"
	})).Return(dirResponse, nil)
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == "https://raw.githubusercontent.com/owner/repo/main/rules/cursor-rules-example/example.mdc"
	})).Return(fileResponse, nil)

	// Create a fetcher
	fetcher := NewGitHubFetcher(mockClient)

	// Call the function being tested
	ctx := context.Background()
	body, err := fetcher.FetchRuleSetByName(ctx, "cursor-rules-example")

	// Assert that there was no error
	assert.NoError(t, err)

	// Read the response body
	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, string(data))

	// Assert that the mock expectations were met
	mockClient.AssertExpectations(t)
}
