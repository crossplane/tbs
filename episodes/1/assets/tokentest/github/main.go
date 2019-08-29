package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "YOUR_TOKEN"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(repos)
}
