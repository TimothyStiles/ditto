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

func Client() *http.Client {
	return &http.Client{
		Transport: &CachingTransport{
			Transport: http.DefaultTransport,
		},
	}
}

type CachingTransport struct {
	Transport http.RoundTripper
}

func (c *CachingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	endpoint := req.URL.String()

	data, err := retrieve(endpoint)
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

	err = cache(endpoint, data)
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader(data))
	return resp, nil
}

func getCacheFilePath(endpoint string) string {
	hash := fnv.New64a()
	hash.Write([]byte(endpoint))
	hashedEndpoint := fmt.Sprintf("%x", hash.Sum(nil))
	return filepath.Join(".ditto", hashedEndpoint)
}

func retrieve(endpoint string) ([]byte, error) {
	cacheFilePath := getCacheFilePath(endpoint)
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil, err
	}
	return os.ReadFile(cacheFilePath)
}

func cache(endpoint string, data []byte) error {
	cacheFilePath := getCacheFilePath(endpoint)
	os.MkdirAll(filepath.Dir(cacheFilePath), os.ModePerm)
	return os.WriteFile(cacheFilePath, data, 0644)
}
