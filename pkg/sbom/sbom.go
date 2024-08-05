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

// GenerateSBOMWithSyft generates an SBOM using the syft binary.
func GenerateSBOMWithSyft(directory, outputFile, repo string) error {
	// Check if syft is installed
	_, err := exec.LookPath("syft")
	if err != nil {
		return fmt.Errorf("syft is not installed or not in PATH: %w", err)
	}
	cmd := exec.Command("syft", "scan", fmt.Sprintf("dir:%s", directory), //nolint:gosec
		"-o", "cyclonedx-json", "--file", outputFile,
		"--select-catalogers", "+github-actions-usage-cataloger", "--source-name", repo)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error generating SBOM with syft: %w\nOutput: %s", err, output)
	}

	return nil
}
