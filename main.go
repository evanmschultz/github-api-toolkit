package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github-api-toolkit/githubfetcher"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	fetcher := githubfetcher.NewRepoFetcher(token)

	fmt.Println("Starting to fetch documentation and examples...")
	err := fetcher.FetchRelevantFiles(context.Background(), "langchain-ai", "langchain")
	if err != nil {
		log.Fatalf("Failed to fetch documentation and examples: %v", err)
	}

	fmt.Println("Documentation and examples fetched successfully")
}