package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bit-bom/bom-factory/pkg/csv"
	"github.com/bit-bom/bom-factory/pkg/sbom"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

const defaultDBPath = "data.db"
const defaultSBOMDir = "sbom"

func main() {
	app := &cli.App{
		Name:  "criticalprojects",
		Usage: "Load CSV data into SQLite and query it",
		Commands: []*cli.Command{
			{
				Name:    "load",
				Aliases: []string{"l"},
				Usage:   "Load CSV data into SQLite",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "db",
						Aliases:  []string{"d"},
						Value:    defaultDBPath,
						Usage:    "Path to the SQLite database file",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "csv",
						Aliases:  []string{"c"},
						Value:    "data.csv",
						Usage:    "Path to the CSV file",
						Required: true,
					},
					&cli.IntFlag{
						Name:  "start",
						Usage: "Start line number (0-based, inclusive)",
					},
					&cli.IntFlag{
						Name:  "end",
						Usage: "End line number (0-based, exclusive, 0 means until the end)",
					},
				},
				Action: loadCSVToSQLite,
			},
			{
				Name:    "query",
				Aliases: []string{"q"},
				Usage:   "Query SQLite data",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "db",
						Aliases:  []string{"d"},
						Value:    defaultDBPath,
						Usage:    "Path to the SQLite database file",
						Required: false,
					},
					&cli.StringSliceFlag{
						Name:     "filter",
						Aliases:  []string{"f"},
						Usage:    "Filter criteria in the format 'field:operator:value' (can be used multiple times)",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "max-results",
						Aliases: []string{"m"},
						Usage:   "Maximum number of results to return",
						Value:   100,
					},
				},
				Action: querySQLiteData,
			},
			{
				Name:    "download-sbom",
				Aliases: []string{"ds"},
				Usage:   "Download SBOM for repositories matching the filter criteria",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "db",
						Aliases:  []string{"d"},
						Value:    defaultDBPath,
						Usage:    "Path to the SQLite database file",
						Required: false,
					},
					&cli.StringSliceFlag{
						Name:     "filter",
						Aliases:  []string{"f"},
						Usage:    "Filter criteria in the format 'field:operator:value' (can be used multiple times)",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "GitHub token for authentication",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "dir",
						Aliases:  []string{"o"},
						Value:    defaultSBOMDir,
						Usage:    "Directory to save the SBOM files",
						Required: false,
					},
					&cli.IntFlag{
						Name:    "max-results",
						Aliases: []string{"m"},
						Usage:   "Maximum number of results to return",
						Value:   100,
					},
				},
				Action: downloadSBOMs,
			},
			{
				Name:    "convert-to-purl",
				Aliases: []string{"cp"},
				Usage:   "Convert SPDX file(s) to include PURLs",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "Path to the SPDX JSON file",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "dir",
						Aliases:  []string{"d"},
						Usage:    "Path to the directory containing SPDX JSON files",
						Required: false,
					},
				},
				Action: convertToPURL,
			},
			{
				Name:    "validate-sbom",
				Aliases: []string{"vs"},
				Usage:   "Validate an SBOM file or all SBOM files in a directory",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "Path to the SBOM file",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "dir",
						Aliases:  []string{"d"},
						Usage:    "Path to the directory containing SBOM files",
						Required: false,
					},
				},
				Action: validateSBOM,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadCSVToSQLite(c *cli.Context) error {
	dbPath := c.String("db")
	csvFilePath := c.String("csv")

	if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
		return fmt.Errorf("CSV file does not exist: %s", csvFilePath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}
	defer db.Close()

	createTableStmt := `
	CREATE TABLE IF NOT EXISTS repos (
		repo_url TEXT PRIMARY KEY,
		repo_language TEXT,
		repo_license TEXT,
		repo_star_count INTEGER,
		repo_created_at TEXT,
		repo_updated_at TEXT,
		legacy_created_since INTEGER,
		legacy_updated_since INTEGER,
		legacy_contributor_count INTEGER,
		legacy_org_count INTEGER,
		legacy_commit_frequency REAL,
		legacy_recent_release_count INTEGER,
		legacy_updated_issues_count INTEGER,
		legacy_closed_issues_count INTEGER,
		legacy_issue_comment_frequency REAL,
		legacy_github_mention_count INTEGER,
		depsdev_dependent_count INTEGER,
		default_score REAL,
		collection_date TEXT,
		worker_commit_id TEXT
	);
	`
	_, err = db.Exec(createTableStmt)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	options := csv.LoadOptions{
		StartLine:  0,
		MaxRecords: 0,
	}

	if c.IsSet("start") {
		options.StartLine = c.Int("start")
	}
	if c.IsSet("end") {
		options.MaxRecords = c.Int("end") - options.StartLine
	}

	err = csv.LoadCSVToSQLite(csvFilePath, db, options)
	if err != nil {
		return fmt.Errorf("failed to load CSV data into SQLite: %w", err)
	}

	if c.IsSet("start") || c.IsSet("end") {
		fmt.Printf("CSV data from %s (lines %d to %d) successfully loaded into SQLite at %s\n",
			csvFilePath, options.StartLine, c.Int("end"), dbPath)
	} else {
		fmt.Printf("CSV data from %s successfully loaded into SQLite at %s\n", csvFilePath, dbPath)
	}
	return nil
}

func querySQLiteData(c *cli.Context) error {
	dbPath := c.String("db")
	filterArgs := c.StringSlice("filter")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}
	defer db.Close()

	var filterCriteria []csv.FilterCriteria
	for _, arg := range filterArgs {
		criterion, err := csv.ParseFilterCriteria(arg)
		if err != nil {
			return fmt.Errorf("invalid filter criteria: %w", err)
		}
		filterCriteria = append(filterCriteria, criterion)
	}

	filteredData, err := csv.FilterSQLiteData(db, filterCriteria, c.Int("max-results")) // Added the third argument
	if err != nil {
		return fmt.Errorf("failed to filter SQLite data: %w", err)
	}

	fmt.Printf("Found %d repositories matching the criteria\n", len(filteredData))

	for i := 0; i < len(filteredData) && i < 5; i++ {
		repo := filteredData[i]
		fmt.Printf("Repo %d: %s (Stars: %d, Language: %s)\n", i+1, repo.RepoURL, repo.RepoStarCount, repo.RepoLanguage)
	}

	return nil
}

func downloadSBOMs(c *cli.Context) error {
	dbPath := c.String("db")
	filterArgs := c.StringSlice("filter")
	token := c.String("token")
	dir := c.String("dir")

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}
	defer db.Close()

	var filterCriteria []csv.FilterCriteria
	for _, arg := range filterArgs {
		criterion, err := csv.ParseFilterCriteria(arg)
		if err != nil {
			return fmt.Errorf("invalid filter criteria: %w", err)
		}
		filterCriteria = append(filterCriteria, criterion)
	}

	filteredData, err := csv.FilterSQLiteData(db, filterCriteria, c.Int("max-results"))
	if err != nil {
		return fmt.Errorf("failed to filter SQLite data: %w", err)
	}

	if len(filteredData) == 0 {
		fmt.Println("No repositories matching the criteria")
		return nil
	}

	// Ensure the directory exists
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	for i := 0; i < len(filteredData); i++ {
		repo := &filteredData[i]
		sbomContent, err := csv.DownloadSBOMFromGitHub(*repo, token) // Dereference the pointer
		if err != nil {
			fmt.Printf("Failed to download SBOM for %s: %v\n", repo.RepoURL, err)
			continue
			// Dereference the pointer when passing it to the function
		}

		parsedURL, err := url.Parse(repo.RepoURL)
		if err != nil {
			fmt.Printf("Failed to parse URL %s: %v\n", repo.RepoURL, err)
			continue
		}

		pathSegments := strings.Split(parsedURL.Path, "/")
		if len(pathSegments) < 3 {
			fmt.Printf("Invalid repository URL format: %s\n", repo.RepoURL)
			continue
		}

		orgName := pathSegments[1]
		repoName := pathSegments[2]
		safeOrgName := url.PathEscape(orgName)
		safeRepoName := url.PathEscape(repoName)
		fileName := fmt.Sprintf("%s_%s", safeOrgName, safeRepoName)

		if err := csv.SaveSBOMToFile(sbomContent, dir, fileName); err != nil {
			fmt.Printf("Failed to save SBOM for %s: %v\n", repo.RepoURL, err)
			continue
		}

		fmt.Printf("SBOM for %s saved successfully\n", repo.RepoURL)

		// Add a delay of 3 seconds between downloads
		time.Sleep(3 * time.Second)
	}

	return nil
}

func convertToPURL(c *cli.Context) error {
	filePath := c.String("file")
	dirPath := c.String("dir")

	if filePath == "" && dirPath == "" {
		return fmt.Errorf("either --file or --dir must be specified")
	}

	if filePath != "" {
		err := sbom.UpdateSPDXWithPURLs(filePath)
		if err != nil {
			return fmt.Errorf("failed to convert SPDX to PURLs: %w", err)
		}
		fmt.Printf("Successfully converted SPDX file to include PURLs: %s\n", filePath)
		return nil
	}

	if dirPath != "" {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}

		var failedFiles []string
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
				filePath := fmt.Sprintf("%s/%s", dirPath, file.Name())
				err := sbom.UpdateSPDXWithPURLs(filePath)
				if err != nil {
					failedFiles = append(failedFiles, fmt.Sprintf("%s: %v", filePath, err))
				}
			}
		}

		if len(failedFiles) > 0 {
			for _, failure := range failedFiles {
				fmt.Println(failure)
			}
			return fmt.Errorf("conversion failed for %d file(s)", len(failedFiles))
		}

		fmt.Printf("Successfully converted all SPDX files in directory: %s\n", dirPath)
		return nil
	}

	return nil
}

func validateSBOM(c *cli.Context) error {
	filePath := c.String("file")
	dirPath := c.String("dir")

	if filePath == "" && dirPath == "" {
		return fmt.Errorf("either --file or --dir must be specified")
	}

	var failedFiles []string

	if filePath != "" {
		err := sbom.ValidateSBOM(filePath)
		if err != nil {
			failedFiles = append(failedFiles, fmt.Sprintf("%s: %v", filePath, err))
		}
	} else if dirPath != "" {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
				filePath := fmt.Sprintf("%s/%s", dirPath, file.Name())
				err := sbom.ValidateSBOM(filePath)
				if err != nil {
					failedFiles = append(failedFiles, fmt.Sprintf("%s: %v", filePath, err))
				}
			}
		}
	}

	if len(failedFiles) > 0 {
		for _, failure := range failedFiles {
			fmt.Println(failure)
		}
		return fmt.Errorf("validation failed for %d file(s)", len(failedFiles))
	}

	fmt.Println("All SBOM files validated successfully")
	return nil
}
