package sbom

import (
	"context"
	"fmt"
	"os/exec"
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
	_, err := exec.LookPath("cdxgen")
	if err != nil {
		return fmt.Errorf("cdxgen is not installed or not in PATH: %w", err)
	}

	// Create a context with a 5-minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cdxgen", "-r", "-o", outputFile, "--no-install-deps", "--project-name", repo, "--install-deps", "false", "--spec-version", "1.5", directory)
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
