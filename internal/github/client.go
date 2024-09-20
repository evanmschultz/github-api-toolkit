package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client struct {
	v4 *githubv4.Client
}

func NewClient(token string) (*Client, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	
	v4Client := githubv4.NewClient(httpClient)

	return &Client{
		v4: v4Client,
	}, nil
}