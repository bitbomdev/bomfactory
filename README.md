
# bomfactory

## Overview

`bomfactory` is a powerful command-line tool that simplifies and automates the process of downloading Software Bill of Materials (SBOMs) for multiple repositories. SBOMs are essential for software testing and security analysis, and `bomfactory` streamlines the task of acquiring them efficiently. 

Inspired by the [criticality_score](https://github.com/ossf/criticality_score) project, `bomfactory` specifically targets critical projects, making it an indispensable tool for anyone involved in software security, testing, or research.

## Why Use bomfactory?

Working with thousands of SBOMs can be time-consuming and complex. We developed `bomfactory` to facilitate the rapid acquisition of SBOMs, ensuring that important projects are well-supported. This tool was crucial in testing thousands of SBOMs for [minefield](https://github.com/bitbomdev/minefield), a simple graph database for managing dependencies. By leveraging Roaring Bitmaps, it allows for O(1) query performance on large datasets, significantly simplifying the process.

If you require a large collection of SBOMs for testing or research, check out the [bom-silo](https://github.com/bitbomdev/bom-silo) repository, which was created using `bomfactory`.

## Key Features

- **Download Criticality Score CSV**: Quickly download a CSV file containing criticality scores.
- **Load CSV into SQLite**: Import CSV data into an SQLite database for efficient querying.
- **Advanced Querying**: Perform complex queries on the SQLite database to identify repositories based on various criteria.
- **Automated SBOM Downloads**: Download SBOMs automatically for repositories that match your query criteria.

## Quickstart

> **Note:** Replace `~/temp` with the path to your preferred directory.

### Step 1: Download the CSV file containing criticality scores

```bash
docker run --rm -v ~/temp:/app/data ghcr.io/bitbomdev/bomfactory download-csv -o /app/data/data.csv
```

### Step 2: Load the CSV data into SQLite

```bash
docker run --rm -v ~/temp:/app/data ghcr.io/bitbomdev/bomfactory load -d /app/data/data.db -c /app/data/data.csv --start 1 --end 1000
```

### Step 3: Query the SQLite data

```bash
docker run --rm -v ~/temp:/app/data ghcr.io/bitbomdev/bomfactory q -d /app/data/data.db -f "repo_language:=:Go"
```

### Step 4: Download SBOMs for repositories

```bash
docker run --rm -v ~/temp:/app/data ghcr.io/bitbomdev/bomfactory ds -d /app/data/data.db -f "repo_language:=:Go" --dir /app/data
```

## Advanced Usage

The following example demonstrates how to download 1,000 SBOMs for Go repositories hosted on Google, skipping the first 9,000 repositories and downloading 10 SBOMs concurrently:

```bash
docker run --rm -v ~/temp:/app/data ghcr.io/bitbomdev/bomfactory ds --filter "repo_language:=:Go" --filter "repo_url:LIKE:%google/%" -m 1000 --dir /app/data/sboms/go -d /app/data/data.db -s 9000 --cd 10
```

> **Tip:** For a complete dataset, ensure that you load the entire CSV data into the SQLite database before performing advanced queries.

## Installation

To install `bomfactory`, clone the repository and build the project:

```bash
git clone https://github.com/bitbomdev/bomfactory.git
cd bomfactory
make build
```

## Detailed Usage

### 1. Download the CSV File

```bash
bomfactory download-csv --url https://www.googleapis.com/download/storage/v1/b/ossf-criticality-score/o/2024.07.05%2F143335%2Fall.csv?generation=1721362287412491&alt=media --output data.csv
```

### 2. Load the CSV Data into SQLite

```bash
bomfactory load --csv data.csv --db data.db --start 1 --end 0
```

### 3. Query the SQLite Data

```bash
bomfactory query --filter "repo_language:==:Go" --filter "repo_star_count:>:100" --db data.db
```

### 4. Download SBOMs for Repositories

```bash
bomfactory download-sbom --filter "repo_language:==:Go" --token my_github_token --dir sbom_files --db data.db
```

## Contributions and Support

We welcome contributions and feedback! If you have any questions or need assistance, feel free to open an issue in the repository.

---

This revised README improves readability and organization, ensuring that users can quickly understand the purpose of `bomfactory` and how to use it effectively. Let me know if there's anything else you'd like to add or modify!