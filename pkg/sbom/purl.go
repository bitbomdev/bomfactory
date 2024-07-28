package sbom

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/package-url/packageurl-go"
)

type SPDX struct {
	SPDXID       string `json:"SPDXID"`
	SpdxVersion  string `json:"spdxVersion"`
	CreationInfo struct {
		Created  string   `json:"created"`
		Creators []string `json:"creators"`
	} `json:"creationInfo"`
	Name              string    `json:"name"`
	DataLicense       string    `json:"dataLicense"`
	DocumentDescribes []string  `json:"documentDescribes"`
	DocumentNamespace string    `json:"documentNamespace"`
	Packages          []Package `json:"packages"`
}

type Package struct {
	SPDXID           string        `json:"SPDXID"`
	Name             string        `json:"name"`
	VersionInfo      string        `json:"versionInfo"`
	DownloadLocation string        `json:"downloadLocation"`
	LicenseDeclared  string        `json:"licenseDeclared"`
	ExternalRefs     []ExternalRef `json:"externalRefs,omitempty"`
	FilesAnalyzed    bool          `json:"filesAnalyzed"`
}

type ExternalRef struct {
	ReferenceCategory string `json:"referenceCategory"`
	ReferenceType     string `json:"referenceType"`
	ReferenceLocator  string `json:"referenceLocator"`
}

func UpdateSPDXWithPURLs(filePath string) error {
	// Read the SPDX JSON file
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Parse the JSON into a Go struct
	var spdx SPDX
	err = json.Unmarshal(file, &spdx)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	// Update the struct to include PURLs
	for i, pkg := range spdx.Packages {
		var purl *packageurl.PackageURL
		nameParts := strings.Split(pkg.Name, ":")
		if len(nameParts) < 2 {
			continue
		}
		pkgType := nameParts[0]
		pkgName := nameParts[1]
		switch pkgType {
		case "go":
			purl = packageurl.NewPackageURL(
				packageurl.TypeGolang,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "actions":
			purl = packageurl.NewPackageURL(
				packageurl.TypeGeneric,
				"github",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "bitbucket":
			purl = packageurl.NewPackageURL(
				packageurl.TypeBitbucket,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "deb":
			purl = packageurl.NewPackageURL(
				packageurl.TypeDebian,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "docker":
			purl = packageurl.NewPackageURL(
				packageurl.TypeDocker,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "gem":
			purl = packageurl.NewPackageURL(
				packageurl.TypeGem,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "github":
			purl = packageurl.NewPackageURL(
				packageurl.TypeGithub,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "maven":
			purl = packageurl.NewPackageURL(
				packageurl.TypeMaven,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "npm":
			purl = packageurl.NewPackageURL(
				packageurl.TypeNPM,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "nuget":
			purl = packageurl.NewPackageURL(
				packageurl.TypeNuget,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "pypi":
			purl = packageurl.NewPackageURL(
				packageurl.TypePyPi,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "rpm":
			purl = packageurl.NewPackageURL(
				packageurl.TypeRPM,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		// Add additional cases for other types
		case "alpm":
			purl = packageurl.NewPackageURL(
				packageurl.TypeAlpm,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "apk":
			purl = packageurl.NewPackageURL(
				packageurl.TypeApk,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "bitnami":
			purl = packageurl.NewPackageURL(
				packageurl.TypeBitnami,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "cargo":
			purl = packageurl.NewPackageURL(
				packageurl.TypeCargo,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "cocoapods":
			purl = packageurl.NewPackageURL(
				packageurl.TypeCocoapods,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "composer":
			purl = packageurl.NewPackageURL(
				packageurl.TypeComposer,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "conan":
			purl = packageurl.NewPackageURL(
				packageurl.TypeConan,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "conda":
			purl = packageurl.NewPackageURL(
				packageurl.TypeConda,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "cran":
			purl = packageurl.NewPackageURL(
				packageurl.TypeCran,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "generic":
			purl = packageurl.NewPackageURL(
				packageurl.TypeGeneric,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "hackage":
			purl = packageurl.NewPackageURL(
				packageurl.TypeHackage,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "hex":
			purl = packageurl.NewPackageURL(
				packageurl.TypeHex,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "huggingface":
			purl = packageurl.NewPackageURL(
				packageurl.TypeHuggingface,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "mlflow":
			purl = packageurl.NewPackageURL(
				packageurl.TypeMLFlow,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "oci":
			purl = packageurl.NewPackageURL(
				packageurl.TypeOCI,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "pub":
			purl = packageurl.NewPackageURL(
				packageurl.TypePub,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "qpkg":
			purl = packageurl.NewPackageURL(
				packageurl.TypeQpkg,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "swid":
			purl = packageurl.NewPackageURL(
				packageurl.TypeSWID,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		case "swift":
			purl = packageurl.NewPackageURL(
				packageurl.TypeSwift,
				"",
				pkgName,
				pkg.VersionInfo,
				nil,
				"",
			)
		default:
			continue
		}

		spdx.Packages[i].ExternalRefs = append(spdx.Packages[i].ExternalRefs, ExternalRef{
			ReferenceCategory: "PACKAGE-MANAGER",
			ReferenceType:     "purl",
			ReferenceLocator:  purl.ToString(),
		})
	}

	// Write the updated struct back to the JSON file
	updatedFile, err := json.MarshalIndent(spdx, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = os.WriteFile(filePath, updatedFile, 0o600)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
