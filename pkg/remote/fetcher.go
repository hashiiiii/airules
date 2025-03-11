package remote

import (
	"context"
	"io"
	"net/http"
	"time"
)

// RuleSet represents a rule set from a remote repository
type RuleSet struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Path        string `json:"path"`
	Type        string `json:"type"` // cursor, windsurf, etc.
}

// Fetcher defines the interface for fetching remote resources
type Fetcher interface {
	// ListRuleSets lists available rule sets from the remote repository
	ListRuleSets(ctx context.Context) ([]RuleSet, error)

	// FetchRuleSet fetches a specific rule set from the remote repository
	FetchRuleSet(ctx context.Context, ruleSet RuleSet) (io.ReadCloser, error)

	// FetchRuleSetByName fetches a rule set by name
	FetchRuleSetByName(ctx context.Context, name string) (io.ReadCloser, error)
}

// HTTPClient is an interface for HTTP clients
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient returns a default HTTP client with reasonable timeouts
func DefaultHTTPClient() HTTPClient {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// FetchURL fetches content from a URL using the provided HTTP client
func FetchURL(ctx context.Context, client HTTPClient, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, &FetchError{
			StatusCode: resp.StatusCode,
			URL:        url,
		}
	}

	return resp.Body, nil
}

// FetchError represents an error that occurred during fetching
type FetchError struct {
	StatusCode int
	URL        string
}

// Error implements the error interface
func (e *FetchError) Error() string {
	return "failed to fetch from " + e.URL + " with status code " + http.StatusText(e.StatusCode)
}
