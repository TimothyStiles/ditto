package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGetCacheFilePath(t *testing.T) {
	endpoint := "https://example.com/api"
	expected := filepath.Join(".cache", "362656336a5f086c")
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

func TestGetCachedResponse(t *testing.T) {
	endpoint := "https://example.com/api"
	expectedData := []byte("test data")
	expectedFilePath := getCacheFilePath(endpoint)

	// Create a cache file with test data
	err := os.WriteFile(expectedFilePath, expectedData, 0644)
	if err != nil {
		t.Fatalf("Failed to create cache file: %v", err)
	}
	defer os.Remove(expectedFilePath)

	// Test loading cached response
	result, err := getCachedResponse(endpoint)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Compare loaded data with expected data
	if !bytes.Equal(result, expectedData) {
		t.Errorf("Expected data: %s, but got: %s", expectedData, result)
	}

	// Clean up the cache file
	err = os.Remove(expectedFilePath)
	if err != nil {
		t.Errorf("Failed to remove cache file: %v", err)
	}
}
func TestMakeAPICall(t *testing.T) {
	endpoint := "https://example.com/api"
	expectedData := []byte("test data")
	expectedFilePath := getCacheFilePath(endpoint)

	// Test making API call
	result, err := makeAPICall(endpoint)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Compare API response with expected data
	if !bytes.Equal(result, expectedData) {
		t.Errorf("Expected data: %s, but got: %s", expectedData, result)
	}

	// Clean up the cache file
	err = os.Remove(expectedFilePath)
	if err != nil {
		t.Errorf("Failed to remove cache file: %v", err)
	}
}
