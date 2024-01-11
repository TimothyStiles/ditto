package ditto

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetCacheDir(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com/api", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	cacheDir := getCacheDir(req)

	expectedCacheDir := ".ditto/316496c36323dd38"
	if cacheDir != expectedCacheDir {
		t.Errorf("Expected cache directory `%s`, but got `%s`", expectedCacheDir, cacheDir)
	}
}

func TestCachingTransport_RoundTrip_CachedResponse(t *testing.T) {
	// Define the test URL and expected response
	url := "https://example.com/api"

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Remove any existing cache for the test URL
	cacheDir := getCacheDir(req)
	err = os.RemoveAll(cacheDir)
	if err != nil {
		t.Fatalf("Failed to remove cache directory: %v", err)
	}

	// Create a new caching HTTP client
	client := &CachingTransport{
		Transport: http.DefaultTransport,
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
	_ = os.Remove(cacheDir)

}
func TestCache(t *testing.T) {
	endpoint := "https://example.com/api"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"message": "Hello, World!"}`)),
	}

	err = cache(req, resp)
	if err != nil {
		t.Fatalf("Failed to cache response: %v", err)
	}

	cacheDir := getCacheDir(req)
	body, err := os.ReadFile(filepath.Join(cacheDir, "body"))
	if err != nil {
		t.Fatalf("Failed to read cached body: %v", err)
	}

	expectedBody := `{"message": "Hello, World!"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body `%s`, but got `%s`", expectedBody, string(body))
	}

	data, err := os.ReadFile(filepath.Join(cacheDir, "data"))
	if err != nil {
		t.Fatalf("Failed to read cached data: %v", err)
	}

	var cachedResp CachedResponse
	err = json.Unmarshal(data, &cachedResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal cached response: %v", err)
	}

	expectedStatusCode := http.StatusOK
	if cachedResp.StatusCode != expectedStatusCode {
		t.Errorf("Expected status code `%d`, but got `%d`", expectedStatusCode, cachedResp.StatusCode)
	}

	expectedContentType := "application/json"
	if cachedResp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Expected Content-Type `%s`, but got `%s`", expectedContentType, cachedResp.Header.Get("Content-Type"))
	}
}
func TestRetrieve(t *testing.T) {
	req, err := http.NewRequest("POST", "https://example.com/api", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Remove any existing cache for the test URL
	cacheDir := getCacheDir(req)
	defer os.RemoveAll(cacheDir)

	// make request and cache it
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"message": "Hello, World!"}`)),
	}

	err = cache(req, resp)
	if err != nil {
		t.Fatalf("Failed to cache response: %v", err)
	}

	// retrieve cached response
	cachedResp, err := retrieve(req)
	if err != nil {
		t.Fatalf("Failed to retrieve cached response: %v", err)
	}

	// check that the response body is the same
	body, err := io.ReadAll(cachedResp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedBody := `{"message": "Hello, World!"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body `%s`, but got `%s`", expectedBody, string(body))
	}
}
