package csv

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var ctx = context.Background()

// LoadOptions defines options for loading CSV data into SQLite
type LoadOptions struct {
	StartLine  int // Line number to start loading from (0-based, excluding header)
	MaxRecords int // Maximum number of records to load (0 means load all)
}

// LoadCSVToSQLite loads CSV data into SQLite
func LoadCSVToSQLite(filePath string, db *sql.DB, options LoadOptions) error {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Map CSV column names to SQLite column names
	columnMap := map[string]string{
		"repo.url":                    "repo_url",
		"repo.language":               "repo_language",
		"repo.license":                "repo_license",
		"repo.star_count":             "repo_star_count",
		"repo.created_at":             "repo_created_at",
		"repo.updated_at":             "repo_updated_at",
		"legacy.created_since":        "legacy_created_since",
		"legacy.updated_since":        "legacy_updated_since",
		"legacy.contributor_count":    "legacy_contributor_count",
		"legacy.org_count":            "legacy_org_count",
		"legacy.commit_frequency":     "legacy_commit_frequency",
		"legacy.recent_release_count": "legacy_recent_release_count",
		"legacy.updated_issues_count": "legacy_updated_issues_count",
		"legacy.closed_issues_count":  "legacy_closed_issues_count",
		"legacy.issue_comment_freq":   "legacy_issue_comment_frequency",
		"legacy.github_mention_count": "legacy_github_mention_count",
		"depsdev.dependent_count":     "depsdev_dependent_count",
		"default_score":               "default_score",
		"collection_date":             "collection_date",
		"worker_commit_id":            "worker_commit_id",
	}

	// Quote and map column names
	for i, col := range header {
		if mappedCol, ok := columnMap[col]; ok {
			header[i] = fmt.Sprintf(`"%s"`, mappedCol)
		} else {
			return fmt.Errorf("unknown column name: %s", col)
		}
	}

	// Skip lines if StartLine is specified
	for i := 0; i < options.StartLine; i++ {
		_, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return nil // Reached end of file while skipping
			}
			return fmt.Errorf("failed to skip to start line: %w", err)
		}
	}

	// Prepare insert statement
	insertStmt := fmt.Sprintf("INSERT INTO repos (%s) VALUES (%s)", strings.Join(header, ","), strings.Repeat("?,", len(header)-1)+"?")
	stmt, err := db.Prepare(insertStmt)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	// Load records into SQLite
	loadedRecords := 0
	for options.MaxRecords == 0 || loadedRecords < options.MaxRecords {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read record: %w", err)
		}

		// Convert record to interface slice and handle empty strings
		values := make([]interface{}, len(record))
		for i, v := range record {
			if v == "" {
				switch header[i] {
				case `"repo_star_count"`, `"legacy_created_since"`, `"legacy_updated_since"`, `"legacy_contributor_count"`, `"legacy_org_count"`, `"legacy_recent_release_count"`, `"legacy_updated_issues_count"`, `"legacy_closed_issues_count"`, `"legacy_github_mention_count"`, `"depsdev_dependent_count"`:
					values[i] = 0
				case `"legacy_commit_frequency"`, `"legacy_issue_comment_frequency"`, `"default_score"`:
					values[i] = 0.0
				default:
					values[i] = nil
				}
			} else {
				values[i] = v
			}
		}

		_, err = stmt.Exec(values...)
		if err != nil {
			return fmt.Errorf("failed to insert record into sqlite: %w", err)
		}
		loadedRecords++
	}

	return nil
}

// Operator represents the comparison operator for filtering
type Operator string

const (
	OperatorEqual              Operator = "="
	OperatorGreaterThan        Operator = ">"
	OperatorLessThan           Operator = "<"
	OperatorGreaterThanOrEqual Operator = ">="
	OperatorLessThanOrEqual    Operator = "<="
	OperatorNotEqual           Operator = "!="
	OperatorLike               Operator = "LIKE"
	OperatorNotLike            Operator = "NOT LIKE"
	OperatorIn                 Operator = "IN"
	OperatorNotIn              Operator = "NOT IN"
)

// FilterCriteria defines the criteria for filtering rows
type FilterCriteria struct {
	Field    string
	Value    string
	Operator Operator
}

// ParseFilterCriteria parses a string into FilterCriteria
func ParseFilterCriteria(criteriaStr string) (FilterCriteria, error) {
	parts := strings.SplitN(criteriaStr, ":", 3)
	if len(parts) != 3 {
		return FilterCriteria{}, fmt.Errorf("invalid filter criteria format: %s", criteriaStr)
	}

	return FilterCriteria{
		Field:    parts[0],
		Operator: Operator(parts[1]),
		Value:    parts[2],
	}, nil
}

// FilterSQLiteData filters data in SQLite based on multiple criteria and returns a slice of RepoData structs
func FilterSQLiteData(db *sql.DB, criteria []FilterCriteria) ([]RepoData, error) {
	var filteredRecords []RepoData

	// Build query
	query := "SELECT * FROM repos WHERE "
	args := []interface{}{}
	for i, criterion := range criteria {
		if i > 0 {
			query += " AND "
		}
		query += fmt.Sprintf("%s %s ?", criterion.Field, criterion.Operator)
		args = append(args, criterion.Value)
	}

	// Add ORDER BY clause for Criticality Score (default_score) in descending order
	query += " ORDER BY default_score DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sqlite: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var repo RepoData
		var depsDevDependentCount sql.NullString
		err := rows.Scan(
			&repo.RepoURL, &repo.RepoLanguage, &repo.RepoLicense, &repo.RepoStarCount,
			&repo.RepoCreatedAt, &repo.RepoUpdatedAt, &repo.LegacyCreatedSince, &repo.LegacyUpdatedSince,
			&repo.LegacyContributorCount, &repo.LegacyOrgCount, &repo.LegacyCommitFrequency, &repo.LegacyRecentReleaseCount,
			&repo.LegacyUpdatedIssuesCount, &repo.LegacyClosedIssuesCount, &repo.LegacyIssueCommentFreq, &repo.LegacyGithubMentionCount,
			&depsDevDependentCount, &repo.DefaultScore, &repo.CollectionDate, &repo.WorkerCommitID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if depsDevDependentCount.Valid && depsDevDependentCount.String != "" {
			count, err := strconv.Atoi(depsDevDependentCount.String)
			if err != nil {
				return nil, fmt.Errorf("failed to convert depsdev_dependent_count to int: %w", err)
			}
			repo.DepsDevDependentCount = count
		} else {
			repo.DepsDevDependentCount = 0
		}
		filteredRecords = append(filteredRecords, repo)
	}

	return filteredRecords, nil
}
