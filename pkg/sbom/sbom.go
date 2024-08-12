package sbom

import (
	"fmt"
	"os/exec"

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
	cmd := exec.Command("cdxgen", "-r", "-o", outputFile, "--install-deps", "false", "--spec-version", "1.5", directory)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error generating SBOM with cdxgen: %w\nOutput: %s", err, output)
	}

	return nil
}
