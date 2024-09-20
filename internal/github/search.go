package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

func (c *Client) SearchTopRepos(language string, limit int) ([]*Repository, error) {
	var query struct {
		Search struct {
			Nodes []struct {
				Repository Repository `graphql:"... on Repository"`
			}
		} `graphql:"search(query: $query, type: REPOSITORY, first: $limit)"`
	}

	variables := map[string]interface{}{
		"query": githubv4.String(fmt.Sprintf("language:%s sort:stars-desc", language)),
		"limit": githubv4.Int(limit),
	}

	err := c.v4.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	repos := make([]*Repository, 0, len(query.Search.Nodes))
	for _, node := range query.Search.Nodes {
		repos = append(repos, &node.Repository)
	}

	return repos, nil
}