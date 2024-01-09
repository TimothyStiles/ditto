package ditto

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGetCacheFilePath(t *testing.T) {
	endpoint := "https://example.com/api"
	expected := filepath.Join(".ditto", "362656336a5f086c")
	result := getCacheFilePath(endpoint)
	if result != expected {
		t.Errorf("Expected` %s, but got %s", expected, result)
	}
}
func TestSaveCache(t *testing.T) {
	endpoint := "https://example.com/api"
	data := []byte("test data")
	expectedFilePath := getCacheFilePath(endpoint)

	err := saveCache(endpoint, data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the cache file exists
	_, err = os.Stat(expectedFilePath)
	if err != nil {
		t.Errorf("Cache file does not exist: %v", err)
	}

	// Clean up the cache file
	err = os.Remove(expectedFilePath)
	if err != nil {
		t.Errorf("Failed to remove cache file: %v", err)
	}
}
func TestLoadCache(t *testing.T) {
	endpoint := "https://example.com/api"
	expectedData := []byte("test data")
	expectedFilePath := getCacheFilePath(endpoint)

	// Create a cache file with test data
	err := os.WriteFile(expectedFilePath, expectedData, 0644)
	if err != nil {
		t.Fatalf("Failed to create cache file: %v", err)
	}
	defer os.Remove(expectedFilePath)

	// Test loading cache
	result, err := loadCache(endpoint)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Compare loaded data with expected data
	if !bytes.Equal(result, expectedData) {
		t.Errorf("Expected data: %s, but got: %s", expectedData, result)
	}
}
func TestCachingHTTPClient_RoundTrip_CachedResponse(t *testing.T) {
	// Define the test URL and expected response
	url := "https://example.com/api"

	// Remove any existing cache for the test URL
	cacheFilePath := getCacheFilePath(url)
	_ = os.Remove(cacheFilePath)

	// Create a new caching HTTP client
	client := &CachingHTTPClient{
		Transport: http.DefaultTransport,
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Make the HTTP request using the caching HTTP client
	resp, err := client.RoundTrip(req)
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// clean up the cache file again just for good measure
	_ = os.Remove(cacheFilePath)

}
