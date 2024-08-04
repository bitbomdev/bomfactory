package csv

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v63/github"
	"golang.org/x/oauth2"
)

// DownloadSBOMFromGitHub downloads the SBOM for a repository from GitHub
func DownloadSBOMFromGitHub(repo RepoData, token string) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Extract owner and repo name from RepoURL
	parts := strings.Split(strings.TrimPrefix(repo.RepoURL, "https://github.com/"), "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid repo URL: %s", repo.RepoURL)
	}
	owner, repoName := parts[0], parts[1]

	sbom, _, err := client.DependencyGraph.GetSBOM(ctx, owner, repoName)
	if err != nil {
		return "", fmt.Errorf("failed to get SBOM from GitHub: %w", err)
	}

	sbomJSON, err := serializeSBOM(sbom)
	if err != nil {
		return "", fmt.Errorf("failed to serialize SBOM to JSON: %w", err)
	}

	return sbomJSON, nil
}

func serializeSBOM(sbom *github.SBOM) (string, error) {
	// Convert *github.SBOM to JSON bytes
	sbomBytes, err := json.Marshal(sbom)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SBOM: %w", err)
	}

	// Unmarshal JSON bytes to map[string]interface{}
	var sbomMap map[string]interface{}
	if err := json.Unmarshal(sbomBytes, &sbomMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal SBOM: %w", err)
	}

	// Extract the inner content of the SBOM
	innerContent, ok := sbomMap["sbom"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid SBOM format")
	}

	// Marshal the inner content to JSON
	innerContentJSON, err := json.Marshal(innerContent)
	if err != nil {
		return "", fmt.Errorf("failed to serialize inner content to JSON: %w", err)
	}

	return string(innerContentJSON), nil
}

// SaveSBOMToFile saves the SBOM content to a file
func SaveSBOMToFile(sbomContent, dir, repoName string) error {
	filePath := filepath.Join(dir, fmt.Sprintf("%s.sbom.json", repoName))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create SBOM file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(sbomContent)
	if err != nil {
		return fmt.Errorf("failed to write SBOM content to file: %w", err)
	}

	return nil
}

// CloneRepo clones a Git repository using HTTP to a specified directory without history
func CloneRepo(repoURL, dir string) error {
	// Convert HTTPS URL to HTTP
	httpURL := strings.Replace(repoURL, "https://", "http://", 1)

	// Clone the repository with depth 1 (shallow clone)
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:   httpURL,
		Depth: 1, // Shallow clone
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}
