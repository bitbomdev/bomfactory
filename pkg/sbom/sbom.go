package sbom

import (
	"fmt"

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
