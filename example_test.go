package ditto_test

import (
	"context"
	"fmt"

	"github.com/TimothyStiles/ditto"
	"github.com/google/go-github/v57/github"
)

func Example_basic() {
	client := github.NewClient(ditto.Client()) // instead of http.DefaultClient we use ditto.Client()

	// Use client...
	repos, _, err := client.Repositories.List(context.Background(), "octocat", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(repos[0].GetName())
	// Output: boysenberry-repo-1

}
