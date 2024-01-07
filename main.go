package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func getCacheFilePath(endpoint string) string {
	hash := fnv.New64a()
	hash.Write([]byte(endpoint))
	hashedEndpoint := fmt.Sprintf("%x", hash.Sum(nil))
	return filepath.Join(".cache", hashedEndpoint)
}

func loadCache(endpoint string) ([]byte, error) {
	cacheFilePath := getCacheFilePath(endpoint)
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil, err
	}
	return os.ReadFile(cacheFilePath)
}

func saveCache(endpoint string, data []byte) error {
	cacheFilePath := getCacheFilePath(endpoint)
	os.MkdirAll(filepath.Dir(cacheFilePath), os.ModePerm)
	return os.WriteFile(cacheFilePath, data, 0644)
}

func getCachedResponse(endpoint string) ([]byte, error) {
	data, err := loadCache(endpoint)
	if err == nil {
		return data, nil
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = saveCache(endpoint, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func makeAPICall(endpoint string) ([]byte, error) {
	return getCachedResponse(endpoint)
}
