package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
)

const (
	// DefaultGitHubAPIURL is the default GitHub API URL
	DefaultGitHubAPIURL = "https://api.github.com"

	// DefaultGitHubRawURL is the default GitHub raw content URL
	DefaultGitHubRawURL = "https://raw.githubusercontent.com"

	// DefaultAwesomeCursorRulesOwner is the owner of the awesome-cursorrules repository
	DefaultAwesomeCursorRulesOwner = "PatrickJS"

	// DefaultAwesomeCursorRulesRepo is the name of the awesome-cursorrules repository
	DefaultAwesomeCursorRulesRepo = "awesome-cursorrules"

	// DefaultAwesomeCursorRulesBranch is the default branch of the awesome-cursorrules repository
	DefaultAwesomeCursorRulesBranch = "main"

	// DefaultAwesomeCursorRulesPath is the path to the rules directory in the awesome-cursorrules repository
	DefaultAwesomeCursorRulesPath = "rules"
)

// GitHubFetcher fetches rule sets from GitHub repositories
type GitHubFetcher struct {
	client    HTTPClient
	apiURL    string
	rawURL    string
	owner     string
	repo      string
	branch    string
	rulesPath string
}

// GitHubContent represents a file or directory in a GitHub repository
type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
	GitURL      string `json:"git_url"`
	HTMLURL     string `json:"html_url"`
}

// NewGitHubFetcher creates a new GitHubFetcher with the provided options
func NewGitHubFetcher(client HTTPClient, opts ...GitHubFetcherOption) *GitHubFetcher {
	if client == nil {
		client = DefaultHTTPClient()
	}

	fetcher := &GitHubFetcher{
		client:    client,
		apiURL:    DefaultGitHubAPIURL,
		rawURL:    DefaultGitHubRawURL,
		owner:     DefaultAwesomeCursorRulesOwner,
		repo:      DefaultAwesomeCursorRulesRepo,
		branch:    DefaultAwesomeCursorRulesBranch,
		rulesPath: DefaultAwesomeCursorRulesPath,
	}

	for _, opt := range opts {
		opt(fetcher)
	}

	return fetcher
}

// GitHubFetcherOption is a function that configures a GitHubFetcher
type GitHubFetcherOption func(*GitHubFetcher)

// WithGitHubAPIURL sets the GitHub API URL
func WithGitHubAPIURL(apiURL string) GitHubFetcherOption {
	return func(f *GitHubFetcher) {
		f.apiURL = apiURL
	}
}

// WithGitHubRawURL sets the GitHub raw content URL
func WithGitHubRawURL(rawURL string) GitHubFetcherOption {
	return func(f *GitHubFetcher) {
		f.rawURL = rawURL
	}
}

// WithGitHubOwner sets the GitHub repository owner
func WithGitHubOwner(owner string) GitHubFetcherOption {
	return func(f *GitHubFetcher) {
		f.owner = owner
	}
}

// WithGitHubRepo sets the GitHub repository name
func WithGitHubRepo(repo string) GitHubFetcherOption {
	return func(f *GitHubFetcher) {
		f.repo = repo
	}
}

// WithGitHubBranch sets the GitHub repository branch
func WithGitHubBranch(branch string) GitHubFetcherOption {
	return func(f *GitHubFetcher) {
		f.branch = branch
	}
}

// WithGitHubRulesPath sets the path to the rules directory in the GitHub repository
func WithGitHubRulesPath(rulesPath string) GitHubFetcherOption {
	return func(f *GitHubFetcher) {
		f.rulesPath = rulesPath
	}
}

// ListRuleSets lists available rule sets from the GitHub repository
func (f *GitHubFetcher) ListRuleSets(ctx context.Context) ([]RuleSet, error) {
	// Construct the URL for the GitHub API request
	apiURL, err := url.Parse(f.apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GitHub API URL: %w", err)
	}

	apiURL.Path = path.Join("repos", f.owner, f.repo, "contents", f.rulesPath)

	// Fetch the contents of the rules directory
	body, err := FetchURL(ctx, f.client, apiURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rule sets from GitHub: %w", err)
	}
	defer body.Close()

	// Parse the JSON response
	var contents []GitHubContent
	if err := json.NewDecoder(body).Decode(&contents); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub API response: %w", err)
	}

	// Convert GitHubContent to RuleSet
	var ruleSets []RuleSet
	for _, content := range contents {
		// Only include directories
		if content.Type != "dir" {
			continue
		}

		// Extract the rule set type from the name
		ruleType := "unknown"
		if strings.Contains(content.Name, "cursor") {
			ruleType = "cursor"
		} else if strings.Contains(content.Name, "windsurf") {
			ruleType = "windsurf"
		}

		// Create a RuleSet
		ruleSet := RuleSet{
			Name:        content.Name,
			Description: "", // GitHub API doesn't provide descriptions for directories
			URL:         content.HTMLURL,
			Path:        content.Path,
			Type:        ruleType,
		}

		ruleSets = append(ruleSets, ruleSet)
	}

	return ruleSets, nil
}

// FetchRuleSet fetches a specific rule set from the GitHub repository
func (f *GitHubFetcher) FetchRuleSet(ctx context.Context, ruleSet RuleSet) (io.ReadCloser, error) {
	// First, we need to find the actual rule file within the rule set directory
	apiURL, err := url.Parse(f.apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GitHub API URL: %w", err)
	}

	apiURL.Path = path.Join("repos", f.owner, f.repo, "contents", ruleSet.Path)

	// Fetch the contents of the rule set directory
	body, err := FetchURL(ctx, f.client, apiURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rule set contents from GitHub: %w", err)
	}
	defer body.Close()

	// Parse the JSON response
	var contents []GitHubContent
	if err := json.NewDecoder(body).Decode(&contents); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub API response: %w", err)
	}

	// Find the rule file (typically a .mdc, .rules, or .cursorrules file)
	var ruleFile *GitHubContent
	for _, content := range contents {
		if content.Type != "file" {
			continue
		}

		// Look for .mdc, .rules, or .cursorrules files
		if strings.HasSuffix(content.Name, ".mdc") ||
			strings.HasSuffix(content.Name, ".rules") ||
			content.Name == ".cursorrules" {
			ruleFile = &content
			break
		}
	}

	if ruleFile == nil {
		return nil, fmt.Errorf("no rule file found in rule set %s", ruleSet.Name)
	}

	// Fetch the raw content of the rule file
	return FetchURL(ctx, f.client, ruleFile.DownloadURL)
}

// FetchRuleSetByName fetches a rule set by name
func (f *GitHubFetcher) FetchRuleSetByName(ctx context.Context, name string) (io.ReadCloser, error) {
	// List all rule sets
	ruleSets, err := f.ListRuleSets(ctx)
	if err != nil {
		return nil, err
	}

	// Find the rule set with the specified name
	var targetRuleSet *RuleSet
	for _, ruleSet := range ruleSets {
		if ruleSet.Name == name {
			targetRuleSet = &ruleSet
			break
		}
	}

	if targetRuleSet == nil {
		return nil, fmt.Errorf("rule set %s not found", name)
	}

	// Fetch the rule set
	return f.FetchRuleSet(ctx, *targetRuleSet)
}
