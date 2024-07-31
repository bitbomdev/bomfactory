
# bomfactory

`bomfactory` is a CLI tool that automates downloading Software Bill of Materials (SBOMs) for multiple repositories. It addresses the lack of tools for easily obtaining SBOMs, which are essential for software testing and security analysis. This project is inspired by the [criticality_score](https://github.com/ossf/criticality_score) project to get SBOMs for critical projects. The `sbom` is downloaded from the github if the project has enabled `dependency graph`.

## Index
- [Overview](#overview)
- [How do I get SBOMs for the top 1000 most critical Go projects with more than 100 stars?](#how-do-i-get-sboms-for-the-top-1000-most-critical-go-projects-with-more-than-100-stars)
  - [1. Download the CSV file](#1-download-the-csv-file)
  - [2. Load the CSV data into SQLite](#2-load-the-csv-data-into-sqlite)
  - [3. Query the SQLite data](#3-query-the-sqlite-data)
  - [4. Download SBOMs for repositories](#4-download-sboms-for-repositories)
- [How can I get top 1000 `go` projects?](#how-can-i-get-top-1000-go-projects)
- [How can I update the SPDX SBOM with `purl`?](#how-can-i-update-the-spdx-sbom-with-purl)
- [Where is the source for the list?](#where-is-the-source-for-the-list)
- [Where is the SBOM being downloaded?](#where-is-the-sbom-being-downloaded)
- [What is the motivation for this project?](#what-is-the-motivation-for-this-project)
- [Commands](#commands)
- [Complex Queries](#complex-queries)
  - [How can I get repositories with more than 500 stars and a specific license?](#how-can-i-get-repositories-with-more-than-500-stars-and-a-specific-license)
  - [How can I get repositories created after a specific date and with a high commit frequency?](#how-can-i-get-repositories-created-after-a-specific-date-and-with-a-high-commit-frequency)
  - [How can I get repositories with a specific language, more than 100 stars, and sorted by star count?](#how-can-i-get-repositories-with-a-specific-language-more-than-100-stars-and-sorted-by-star-count)
  - [How can I get repositories with a specific keyword in their URL and a high default score?](#how-can-i-get-repositories-with-a-specific-keyword-in-their-url-and-a-high-default-score)
  - [How can I get repositories with a high number of dependent projects and a recent release?](#how-can-i-get-repositories-with-a-high-number-of-dependent-projects-and-a-recent-release)
  - [How can I get repositories with a specific organization and a high number of contributors?](#how-can-i-get-repositories-with-a-specific-organization-and-a-high-number-of-contributors)
  - [How can I get repositories with a specific language, more than 100 stars, and a specific license?](#how-can-i-get-repositories-with-a-specific-language-more-than-100-stars-and-a-specific-license)

## Overview
`bomfactory` is a CLI tool to load CSV data into SQLite, query the data, download SBOMs for repositories, and convert SPDX files to include PURLs. This project is inspired by the [criticality_score](https://github.com/ossf/criticality_score) project to get SBOMs for critical projects.

## How do I get SBOMs for the top 1000 most critical Go projects with more than 100 stars?
To get SBOMs for top 1000 most critical Go projects, follow these steps:

#### 1. Download the CSV file
To download the Criticality Score CSV file, use the `download-csv` command:

_URL_ is optional, if not provided, the default URL will be used.

_OUTPUT_ is optional, if not provided, the default output file will be used.

```sh
bomfactory download-csv --url https://www.googleapis.com/download/storage/v1/b/ossf-criticality-score/o/2024.07.05%2F143335%2Fall.csv?generation=1721362287412491&alt=media --output data.csv
```

This will download the CSV file and save it to `data.csv`.

#### 2. Load the CSV data into SQLite
To load the CSV data into an SQLite database, use the `load` command:

_CSV_ is optional, if not provided, the default CSV file will be used.

_DB_ is optional, if not provided, the default database file will be used.

_END_ if 0, it will load all the data.

_START_ is the start line of the data to load.

```sh
bomfactory load --csv data.csv --db data.db --start 1 --end 0
```

This will load the data from line 10 (inclusive) to line 100 (exclusive) into the SQLite database located at `data.db`.

#### 3. Query the SQLite data
To query the data in the SQLite database, use the `query` command:

```sh
bomfactory query --filter "repo_language:==:Go" --filter "repo_star_count:>:100" --db data.db
```

This will return a list of repositories written in Go with more than 100 stars.

#### 4. Download SBOMs for repositories
To download the SBOMs for repositories matching the filter criteria, use the `download-sbom` command:

```sh
bomfactory download-sbom --filter "repo_language:==:Go" --token my_github_token --dir sbom_files --db data.db
```

This will download the SBOMs for all Go repositories and save them in the `sbom_files` directory. You'll need to provide a valid GitHub token for authentication.

### How can I get top 1000 `javascript` projects?
To get the top 1000 javascript projects, you can adjust the `query` command to fetch the top 1000 projects based on your criteria:

```sh
bomfactory query --filter "repo_language:LIKE:%javascript%" -m 1000 
```

### How can I update the SPDX SBOM with `purl`?
If you have SPDX JSON files that you want to update to include PURLs, use the `convert-to-purl` command:

```sh
bomfactory convert-to-purl --file spdx.json
```

This will update the `spdx.json` file to include the PURLs for the packages listed in the SPDX file.

### Where is the source for the list?
The source for the list of repositories is the Criticality Score CSV file, which can be downloaded using the `download-csv` command.

### Where is the SBOM being downloaded?
The SBOMs are downloaded to the directory specified by the `--dir` option in the `download-sbom` command. For example, in the command:

```sh
bomfactory download-sbom --filter "repo_language:==:Go" --token my_github_token --dir sbom_files --db data.db
```

The SBOMs will be saved in the `sbom_files` directory.

### What is the motivation for this project?
The motivation for this project is to provide an easy-to-use tool for obtaining SBOMs, which are essential for software testing and security analysis. The project is inspired by the [criticality_score](https://github.com/ossf/criticality_score) project to get SBOMs for critical projects.

## Commands

The `bomfactory` tool supports the following commands:

- `download-csv`: Download the Criticality Score CSV file.
- `load`: Load CSV data into SQLite.
- `query`: Query SQLite data.
- `download-sbom`: Download SBOM for repositories matching the filter criteria.
- `convert-to-purl`: Convert SPDX file to include PURLs.

## Complex Queries

### How can I get repositories with more than 500 stars and a specific license?
To get repositories with more than 500 stars and a specific license (e.g., MIT), use the following query:

```sh
bomfactory query --filter "repo_star_count:>:500" --filter "repo_license:==:MIT" --db data.db
```

### How can I get repositories created after a specific date and with a high commit frequency?
To get repositories created after January 1, 2020, and with a commit frequency greater than 10, use the following query:

```sh
bomfactory query --filter "repo_created_at:>:2020-01-01" --filter "legacy_commit_frequency:>:10" --db data.db
```

### How can I get repositories with a specific language, more than 100 stars ?
To get repositories written in Python with more than 100 stars, use the following query:

```sh
bomfactory query --filter "repo_language:==:Python" --filter "repo_star_count:>:100" --db data.db 
```

### How can I get repositories with a specific keyword in their URL and a high default score?
To get repositories with "github" in their URL and a default score greater than 0.8, use the following query:

```sh
bomfactory query --filter "repo_url:LIKE:%github%" --filter "default_score:>:0.8" --db data.db
```

### How can I get repositories with a high number of dependent projects and a recent release?
To get repositories with more than 50 dependent projects and at least one release in the last year, use the following query:

```sh
bomfactory query --filter "depsdev_dependent_count:>:50" --filter "legacy_recent_release_count:>:0" --db data.db
```

### How can I get repositories with a specific organization and a high number of contributors?
To get repositories belonging to the "google" organization with more than 10 contributors, use the following query:

```sh
bomfactory query --filter "repo_url:LIKE:%github.com/google/%" --filter "legacy_contributor_count:>:10" --db data.db
```

### How can I get repositories with a specific language, more than 100 stars, and a specific license?
To get repositories written in JavaScript with more than 100 stars and an Apache-2.0 license, use the following query:

```sh
bomfactory query --filter "repo_language:==:JavaScript" --filter "repo_star_count:>:100" --filter "repo_license:==:Apache-2.0" --db data.db
```
```