package main

import (
	"fmt"
	"log"
	"os"

	"github-api-toolkit/internal/github"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	client, err := github.NewClient(token)
	if err != nil {
		log.Fatalf("Error creating GitHub client: %v", err)
	}

	// Example: Fetch repository information
	repo, err := client.GetRepo("octocat", "Hello-World")
	if err != nil {
		log.Fatalf("Error fetching repository: %v", err)
	}

	fmt.Printf("Repository: %s\nStars: %d\nDescription: %s\n", 
		repo.FullName, repo.StargazersCount, repo.Description)
}