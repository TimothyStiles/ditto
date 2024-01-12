package ditto

import (
	"bytes"
	"encoding/json"
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

type CachedResponse struct {
	StatusCode int
	Status     string
	Method     string
	URL        string
	Header     http.Header
	Body       string
}

// RoundTrip implements the RoundTripper interface.
func (c *CachingTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// So I originally tried implementing this to be a lot cleaner with a majority of the logic in retrieve and cache functions. However,
	// I'm not sure what was going on but there was some sort of data race where the cache file was being written to correctly but if
	// RoundTrip called cache then the GitHub example test would panic and not receive the data from the same response the cache file was
	// being written from in the same call to this function (RoundTrip).

	data, err := retrieve(req)
	if err == nil {

		var cachedResp CachedResponse
		err = json.Unmarshal(data, &cachedResp)
		if err != nil {
			return nil, err
		}
		reader := io.NopCloser(bytes.NewReader([]byte(cachedResp.Body)))
		return &http.Response{
			StatusCode: cachedResp.StatusCode,
			Status:     cachedResp.Status,
			Header:     cachedResp.Header,
			Body:       reader,
		}, nil
	}

	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cachedResp := CachedResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		URL:        req.URL.String(),
		Method:     req.Method,
		Header:     resp.Header,
		Body:       string(body),
	}

	marshalledResponse, err := json.Marshal(cachedResp)
	if err != nil {
		return nil, err
	}

	err = cache(req, marshalledResponse)
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader(body))
	return resp, nil
}

func getCacheFilePath(req *http.Request) string {
	hash := fnv.New64a()
	url := req.URL.String()
	method := req.Method
	endpointPlusMethod := fmt.Sprintf("%s:%s", method, url)
	hash.Write([]byte(endpointPlusMethod))
	hashedEndpoint := fmt.Sprintf("%x", hash.Sum(nil))

	goModDir, err := findGoModDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(goModDir, ".ditto", hashedEndpoint)
}

func retrieve(req *http.Request) ([]byte, error) {
	cacheFilePath := getCacheFilePath(req)
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil, err
	}
	return os.ReadFile(cacheFilePath)
}

func cache(req *http.Request, data []byte) error {
	cacheFilePath := getCacheFilePath(req)
	os.MkdirAll(filepath.Dir(cacheFilePath), os.ModePerm)
	return os.WriteFile(cacheFilePath, data, 0644)
}

func findGoModDir() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	for {
		if path == home || path == "//" {
			return "", nil
		}

		if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
			return path, nil
		}

		path = filepath.Dir(path)
	}
}
