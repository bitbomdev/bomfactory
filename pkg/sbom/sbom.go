package sbom

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	proto "github.com/protobom/protobom/pkg/reader"
)

// ValidateSBOM validates the SBOM file.
func ValidateSBOM(sbom string) error {
	sbomReader := proto.New()
	_, err := sbomReader.ParseFile(sbom)
	if err != nil {
		return fmt.Errorf("error parsing SBOM: %w", err)
	}
	return nil
}

// GenerateSBOMWithCycloneDX generates an SBOM using the cdxgen binary.
func GenerateSBOMWithCycloneDX(directory, outputFile, repo string) error {
	// Check if cdxgen is installed
	_, err := exec.LookPath("syft")
	if err != nil {
		return fmt.Errorf("syft is not installed or not in PATH: %w", err)
	}
	// Generate the output file name by replacing slashes with underscores and appending .json extension
	escapedRepo := strings.ReplaceAll(repo, "/", "_")
	outputFileName := fmt.Sprintf("%s.json", escapedRepo)
	// Create a context with a 5-minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "syft", "scan", directory, "-o", "cyclonedx-json@1.5", "--file", outputFileName)
	fmt.Println("Executing command: for the repo", repo, cmd.String())
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timed out")
	}
	if err != nil {
		return fmt.Errorf("error generating SBOM with cdxgen: %w\nOutput: %s", err, output)
	}

	return nil
}
