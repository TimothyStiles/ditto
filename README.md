# ditto

Ditto is a dead simple code-gen free, mock free, CLI tool free, API call caching package for testing 3rd party APIs in Go.

You shouldn't have to mock 3rd party APIs to test your code. Ditto is a simple package that caches API responses for you to use in your tests.

All you have to do is replace your API client's `http.Client` with `ditto.Client` when writing your tests and you're good to go. Ditto checks if the request has been made before and if so, returns the cached response. If not, it makes the request and caches the response for you to run your tests against later on.


## Install 

```bash
go get github.com/TimothyStiles/ditto
```

## Usage
https://github.com/TimothyStiles/ditto/blob/80c31ef9b298933da6e1a5ac03ce5e189e04ed85/example_test.go#L12-L29
