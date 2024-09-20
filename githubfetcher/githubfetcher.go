package githubfetcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// RepoFetcher handles fetching and saving GitHub repository contents
type RepoFetcher struct {
	client *githubv4.Client
}

// NewRepoFetcher creates a new RepoFetcher instance
func NewRepoFetcher(token string) *RepoFetcher {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	return &RepoFetcher{client: client}
}

// FetchRelevantFiles fetches relevant documentation and examples from a GitHub repository
func (rf *RepoFetcher) FetchRelevantFiles(ctx context.Context, owner, name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	repoDir := filepath.Join(cwd, name+"-relevant-files")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", repoDir, err)
	}

	log.Printf("Fetching relevant files from repository %s/%s into %s", owner, name, repoDir)
	fmt.Printf("%s\n", name+"-relevant-files")

	return rf.fetchRelevantTree(ctx, owner, name, "", repoDir, 0, true)
}

func (rf *RepoFetcher) fetchRelevantTree(ctx context.Context, owner, name, path, localPath string, depth int, isRootDir bool) error {
	// Remove depth limit, but add a safeguard for extremely large repos
	if depth > 10 {
		log.Printf("Warning: Reached maximum depth of 10 at path %s. Skipping deeper levels.", path)
		return nil
	}

	var query struct {
		Repository struct {
			Object struct {
				Tree struct {
					Entries []struct {
						Name string
						Type string
						Object struct {
							Blob struct {
								Text string
							} `graphql:"... on Blob"`
						}
					}
				} `graphql:"... on Tree"`
			} `graphql:"object(expression: $ref)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
		"ref":   githubv4.String("HEAD:" + path),
	}

	err := rf.client.Query(ctx, &query, variables)
	if err != nil {
		return fmt.Errorf("failed to query repository contents for path %s: %w", path, err)
	}

	for _, entry := range query.Repository.Object.Tree.Entries {
		entryPath := filepath.Join(localPath, entry.Name)
		relPath := filepath.Join(path, entry.Name)

		if entry.Type == "tree" {
			if isRelevantDir(entry.Name, isRootDir) {
				fmt.Printf("%s%s%s\n", strings.Repeat("│   ", depth), "├── ", entry.Name)
				if err := os.MkdirAll(entryPath, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", entryPath, err)
				}
				if err := rf.fetchRelevantTree(ctx, owner, name, relPath, entryPath, depth+1, false); err != nil {
					return err
				}
			}
		} else if entry.Type == "blob" && isRelevantFile(entry.Name, isRootDir, path) {
			fmt.Printf("%s%s%s\n", strings.Repeat("│   ", depth), "├── ", entry.Name)
			if err := os.WriteFile(entryPath, []byte(entry.Object.Blob.Text), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", entryPath, err)
			}
		}
	}

	return nil
}

func isRelevantDir(name string, isRootDir bool) bool {
	lowerName := strings.ToLower(name)
	
	relevantDirs := []string{"doc", "docs", "example", "examples", "tutorial", "tutorials", "guide", "guides", "src", "lib", "cookbook"}
	
	if isRootDir {
		relevantRootDirs := append(relevantDirs, "test", "tests")
		for _, dir := range relevantRootDirs {
			if lowerName == dir {
				return true
			}
		}
		return false
	}

	for _, dir := range relevantDirs {
		if strings.Contains(lowerName, dir) {
			return true
		}
	}
	return false
}

func isRelevantFile(name string, isRootDir bool, path string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	lowerName := strings.ToLower(name)

	// Always include certain files in the root directory
	if isRootDir {
		importantRootFiles := []string{"readme", "changelog", "contributing", "authors", "license", "security", "getting_started", "quickstart", "install", "setup"}
		for _, file := range importantRootFiles {
			if strings.Contains(lowerName, file) {
				return true
			}
		}
	}

	// Check if the file is in a known important directory
	pathLower := strings.ToLower(path)
	isInRelevantDir := strings.Contains(pathLower, "doc") || 
		strings.Contains(pathLower, "example") || 
		strings.Contains(pathLower, "tutorial") || 
		strings.Contains(pathLower, "cookbook") || 
		strings.Contains(pathLower, "src") || 
		strings.Contains(pathLower, "lib")

	// Documentation files
	if ext == ".md" || ext == ".mdx" || ext == ".rst" || (ext == ".txt" && strings.Contains(lowerName, "readme")) {
		return true
	}

	// Jupyter notebooks
	if ext == ".ipynb" {
		return true
	}

	// Python files in relevant directories
	if ext == ".py" && isInRelevantDir {
		return true
	}

	// Include JSON files in relevant directories (often used for configuration or examples)
	if ext == ".json" && isInRelevantDir {
		return true
	}

	// Explicitly exclude certain files
	excludedFiles := []string{".gitignore", "requirements.txt", "setup.py", "Makefile", "docker-compose.yml"}
	for _, excluded := range excludedFiles {
		if lowerName == excluded {
			return false
		}
	}

	return false
}