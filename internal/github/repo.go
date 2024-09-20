package github

import (
	"context"

	"github.com/shurcooL/githubv4"
)

type Repository struct {
	ID              string
	Name            string
	FullName        string
	Description     string
	StargazersCount int
}

func (c *Client) GetRepo(owner, name string) (*Repository, error) {
	var query struct {
		Repository struct {
			ID              githubv4.ID
			Name            githubv4.String
			NameWithOwner   githubv4.String
			Description     githubv4.String
			StargazerCount  githubv4.Int
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	err := c.v4.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	return &Repository{
		ID:              string(query.Repository.ID.(githubv4.String)),
		Name:            string(query.Repository.Name),
		FullName:        string(query.Repository.NameWithOwner),
		Description:     string(query.Repository.Description),
		StargazersCount: int(query.Repository.StargazerCount),
	}, nil
}