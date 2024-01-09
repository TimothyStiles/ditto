package ditto

import (
	"bytes"
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

type CachingHTTPClient struct {
	Transport http.RoundTripper
}

func (c *CachingHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	endpoint := req.URL.String()

	data, err := loadCache(endpoint)
	if err == nil {
		reader := io.NopCloser(bytes.NewReader(data))
		return &http.Response{
			StatusCode: 200,
			Body:       reader,
		}, nil
	}

	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = saveCache(endpoint, data)
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader(data))
	return resp, nil
}
