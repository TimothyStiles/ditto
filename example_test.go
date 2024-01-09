package ditto_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TimothyStiles/ditto"
	"github.com/google/go-github/v57/github"
)

func Example_basic() {
	client := github.NewClient(&http.Client{
		Transport: &ditto.CachingHTTPClient{
			Transport: http.DefaultTransport,
		},
	})

	// Use client...
	repos, _, err := client.Repositories.List(context.Background(), "octocat", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(repos[0].GetName())
	// Output: boysenberry-repo-1

}
