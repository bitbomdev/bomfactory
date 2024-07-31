# bomfactory

`bomfactory` is a CLI tool that automates downloading Software Bill of Materials (SBOMs) for multiple repositories. It addresses the lack of tools for easily obtaining SBOMs, which are essential for software testing and security analysis. Inspired by the [criticality_score](https://github.com/ossf/criticality_score) project, `bomfactory` downloads SBOMs from GitHub if the project has enabled the `dependency graph`.

## Table of Contents
- [Overview](#overview)
- [Getting Started](#getting-started)
  - [1. Download the CSV File](#1-download-the-csv-file)
  - [2. Load the CSV Data into SQLite](#2-load-the-csv-data-into-sqlite)
  - [3. Query the SQLite Data](#3-query-the-sqlite-data)
  - [4. Download SBOMs for Repositories](#4-download-sboms-for-repositories)
- [Common Tasks](#common-tasks)
  - [Get Top 1000 Go Projects](#get-top-1000-go-projects)
  - [Update SPDX SBOM with PURLs](#update-spdx-sbom-with-purls)
  - [Source of the List](#source-of-the-list)
  - [SBOM Download Location](#sbom-download-location)
  - [Project Motivation](#project-motivation)
- [Commands](#commands)
- [Complex Queries](#complex-queries)

## Overview
`bomfactory` is a CLI tool for loading CSV data into SQLite, querying the data, downloading SBOMs for repositories, and converting SPDX files to include PURLs. Inspired by the [criticality_score](https://github.com/ossf/criticality_score) project, it targets obtaining SBOMs for critical projects.

## Getting Started
To get SBOMs for the top 1000 most critical Go projects with more than 100 stars, follow these steps:

### 1. Download the CSV File
To download the Criticality Score CSV file, use the `download-csv` command. The URL and output file are optional.

```sh
bomfactory download-csv --url <optional_url> --output data.csv
```

This downloads the CSV file and saves it to `data.csv`.

### 2. Load the CSV Data into SQLite
To load the CSV data into an SQLite database, use the `load` command. CSV and DB file locations are optional.

```sh
bomfactory load --csv data.csv --db data.db --start 1 --end 0
```

This loads the data into the SQLite database at `data.db`.

### 3. Query the SQLite Data
To query the SQLite database, use the `query` command:

```sh
bomfactory query --filter "repo_language:==:Go" --filter "repo_star_count:>:100" --db data.db
```

This returns a list of Go repositories with more than 100 stars.

### 4. Download SBOMs for Repositories
To download SBOMs for repositories matching the filter criteria, use the `download-sbom` command:

```sh
bomfactory download-sbom --filter "repo_language:==:Go" --token <github_token> --dir sbom_files --db data.db
```

This downloads the SBOMs and saves them in the `sbom_files` directory.

## Common Tasks

### Get Top 1000 Go Projects
To get the top 1000 Go projects, adjust the `query` command:

```sh
bomfactory query --filter "repo_language:LIKE:%Go%" -m 1000 
```

### Update SPDX SBOM with PURLs
To update SPDX JSON files to include PURLs, use the `convert-to-purl` command:

```sh
bomfactory convert-to-purl --file spdx.json
```

### Source of the List
The source for the list of repositories is the Criticality Score CSV file, downloadable with the `download-csv` command.

### SBOM Download Location
SBOMs are downloaded to the directory specified by the `--dir` option in the `download-sbom` command.

### Project Motivation
This project aims to provide an easy-to-use tool for obtaining SBOMs, essential for software testing and security analysis. Inspired by the [criticality_score](https://github.com/ossf/criticality_score) project, it targets critical projects.

## Commands
`bomfactory` supports the following commands:
- `download-csv`: Download the Criticality Score CSV file.
- `load`: Load CSV data into SQLite.
- `query`: Query SQLite data.
- `download-sbom`: Download SBOMs for repositories matching filter criteria.
- `convert-to-purl`: Convert SPDX files to include PURLs.

## Complex Queries

### More than 500 Stars and Specific License
```sh
bomfactory query --filter "repo_star_count:>:500" --filter "repo_license:==:MIT" --db data.db
```

### Created After Specific Date with High Commit Frequency
```sh
bomfactory query --filter "repo_created_at:>:2020-01-01" --filter "legacy_commit_frequency:>:10" --db data.db
```

### Specific Language, More than 100 Stars, Sorted by Star Count
```sh
bomfactory query --filter "repo_language:==:Python" --filter "repo_star_count:>:100" --db data.db --order-by "repo_star_count DESC"
```

### Specific Keyword in URL with High Default Score
```sh
bomfactory query --filter "repo_url:LIKE:%github%" --filter "default_score:>:0.8" --db data.db
```

### High Number of Dependent Projects and Recent Release
```sh
bomfactory query --filter "depsdev_dependent_count:>:50" --filter "legacy_recent_release_count:>:0" --db data.db
```

### Specific Organization with High Number of Contributors
```sh
bomfactory query --filter "repo_url:LIKE:%github.com/google/%" --filter "legacy_contributor_count:>:10" --db data.db
```

### Specific Language, More than 100 Stars, Specific License
```sh
bomfactory query --filter "repo_language:==:JavaScript" --filter "repo_star_count:>:100" --filter "repo_license:==:Apache-2.0" --db data.db
```

