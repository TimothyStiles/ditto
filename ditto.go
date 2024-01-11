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

func (c *CachingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	data, err := retrieve(req)
	if err == nil {
		return data, nil
	}

	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	err = cache(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getCacheDir(req *http.Request) string {
	hash := fnv.New64a()
	// add endpoint and method together to get a unique hash
	endpoint := req.URL.String()
	method := req.Method
	stringToHash := fmt.Sprintf("%s%s", endpoint, method)

	hash.Write([]byte(stringToHash))
	hashedEndpoint := fmt.Sprintf("%x", hash.Sum(nil))
	return filepath.Join(".ditto", hashedEndpoint)
}

func retrieve(req *http.Request) (*http.Response, error) {
	cacheDir := getCacheDir(req)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return nil, err
	}

	data, err := os.ReadFile(cacheDir)
	if err != nil {
		return nil, err
	}

	cachedResp := CachedResponse{}
	err = json.Unmarshal(data, &cachedResp)
	if err != nil {
		return nil, err
	}

	bodyBytes := []byte(cachedResp.Body) // Convert string to byte slice
	reader := io.NopCloser(bytes.NewReader(bodyBytes))
	return &http.Response{
		StatusCode: cachedResp.StatusCode,
		Status:     cachedResp.Status,
		Header:     cachedResp.Header,
		Body:       reader,
	}, nil
}

func cache(req *http.Request, resp *http.Response) error {
	cacheDir := getCacheDir(req)
	os.MkdirAll(filepath.Dir(cacheDir), os.ModePerm)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cachedResp := CachedResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		URL:        req.URL.String(),
		Method:     req.Method,
		Header:     resp.Header,
		Body:       string(body),
	}

	data, err := json.Marshal(cachedResp)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheDir, data, 0644)
}
