# bom-factory

## Overview
`bom-factory` is a CLI tool to load CSV data into SQLite, query the data, download SBOMs for repositories, and convert SPDX files to include PURLs. This project is inspired by the [criticality_score](https://github.com/ossf/criticality_score) project to get SBOMs for critical projects.


## Usage
### Load CSV Data into SQLite
Load CSV data into an SQLite database.

```sh
bom-factory load --csv <path_to_csv> [--db <path_to_db>] [--start <start_line>] [--end <end_line>]
```

### Query SQLite Data
Query data from the SQLite database with filter criteria.

```sh
bom-factory query --filter <field:operator:value> [--db <path_to_db>]
```

### Download SBOMs
Download SBOMs for repositories matching the filter criteria.

```sh
bom-factory download-sbom --filter <field:operator:value> --token <github_token> [--db <path_to_db>] [--dir <output_directory>]
```

### Convert SPDX to PURLs
Convert an SPDX JSON file to include PURLs.

```sh
bom-factory convert-to-purl --file <path_to_spdx_json>
```

## Commands
- `load`: Load CSV data into SQLite.
  - `--csv, -c`: Path to the CSV file (required).
  - `--db, -d`: Path to the SQLite database file (default: `data.db`).
  - `--start`: Start line number (0-based, inclusive).
  - `--end`: End line number (0-based, exclusive, 0 means until the end).

- `query`: Query SQLite data.
  - `--filter, -f`: Filter criteria in the format `field:operator:value` (required, can be used multiple times).
  - `--db, -d`: Path to the SQLite database file (default: `data.db`).

- `download-sbom`: Download SBOM for repositories matching the filter criteria.
  - `--filter, -f`: Filter criteria in the format `field:operator:value` (required, can be used multiple times).
  - `--token, -t`: GitHub token for authentication (required).
  - `--db, -d`: Path to the SQLite database file (default: `data.db`).
  - `--dir, -o`: Directory to save the SBOM files (default: `sbom`).

- `convert-to-purl`: Convert SPDX file to include PURLs.
  - `--file, -f`: Path to the SPDX JSON file (required).

## Examples
### Load CSV Data
```sh
bom-factory load --csv data.csv --db mydatabase.db --start 10 --end 100
```

### Query Data
```sh
bom-factory query --filter "repo_language:==:Go" --filter "repo_star_count:>:100"
```

### Download SBOMs
```sh
bom-factory download-sbom --filter "repo_language:==:Go" --token my_github_token --dir sbom_files
```

### Convert SPDX to PURLs
```sh
bom-factory convert-to-purl --file spdx.json
```
