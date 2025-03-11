package remote

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock implementation of HTTPClient
type MockHTTPClient struct {
	mock.Mock
}

// Do implements the HTTPClient interface
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// TestFetchURL tests the FetchURL function
func TestFetchURL(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)

	// Set up the mock response
	responseBody := "test response"
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
	}

	// Set up the mock expectation
	mockClient.On("Do", mock.Anything).Return(response, nil)

	// Call the function being tested
	ctx := context.Background()
	body, err := FetchURL(ctx, mockClient, "https://example.com")

	// Assert that there was no error
	assert.NoError(t, err)

	// Read the response body
	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.Equal(t, responseBody, string(data))

	// Assert that the mock expectations were met
	mockClient.AssertExpectations(t)
}

// TestFetchURLError tests the FetchURL function when an error occurs
func TestFetchURLError(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)

	// Set up the mock response
	response := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString("not found")),
	}

	// Set up the mock expectation
	mockClient.On("Do", mock.Anything).Return(response, nil)

	// Call the function being tested
	ctx := context.Background()
	_, err := FetchURL(ctx, mockClient, "https://example.com")

	// Assert that there was an error
	assert.Error(t, err)
	assert.IsType(t, &FetchError{}, err)

	// Assert that the mock expectations were met
	mockClient.AssertExpectations(t)
}

// TestDefaultHTTPClient tests the DefaultHTTPClient function
func TestDefaultHTTPClient(t *testing.T) {
	// Call the function being tested
	client := DefaultHTTPClient()

	// Assert that the client is not nil
	assert.NotNil(t, client)
}
