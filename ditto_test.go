package ditto

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCachingTransport_RoundTrip(t *testing.T) {
	// Create a test server
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))

	// Set the server to listen on a specific address
	server.Listener, _ = net.Listen("tcp", "localhost:56560")

	server.Start()
	defer server.Close()

	// Create a new CachingTransport
	cachingTransport := &CachingTransport{
		Transport: http.DefaultTransport,
	}

	// Create a new request
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// remove the cache file if it exists and defer removing it again until the end of the test
	cacheFilePath := getCacheFilePath(req)
	os.Remove(cacheFilePath)
	defer os.Remove(cacheFilePath)

	// Call the RoundTrip method
	resp, err := cachingTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// TODO: Add more assertions as needed
}
func TestFindGoModDir(t *testing.T) {
	_, err := findGoModDir()
	if err != nil {
		t.Fatal(err)
	}
}
