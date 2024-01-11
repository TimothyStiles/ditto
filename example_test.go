package ditto_test

import (
	"context"
	"fmt"
	"os"

	"github.com/TimothyStiles/ditto"
	"github.com/google/go-github/v57/github"
)

func Example_basic() {
	token := os.Getenv("GITHUB_TOKEN")

	// instead of http.DefaultClient we use ditto.Client()
	client := github.NewClient(ditto.Client()).WithAuthToken(token) // auth token is optional

	// Use client...
	repos, _, err := client.Repositories.List(context.Background(), "octocat", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(repos[0].GetName())
	// Output: boysenberry-repo-1

}
