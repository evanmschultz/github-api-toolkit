package github

import (
	"context"
	"encoding/json"
)

func (c *Client) FetchSchema() (string, error) {
	var query struct {
		Schema struct {
			Types []struct {
				Name        string
				Description string
			}
		} `graphql:"__schema"`
	}

	err := c.v4.Query(context.Background(), &query, nil)
	if err != nil {
		return "", err
	}

	schemaJSON, err := json.MarshalIndent(query.Schema, "", "  ")
	if err != nil {
		return "", err
	}

	return string(schemaJSON), nil
}