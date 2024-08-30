# bomfactory

## Why?

`bomfactory` makes it easy to get SBOMs quickly, saving time and ensuring important projects. We aimed to test thousands of SBOMs for [minefield](https://github.com/bitbomdev/minefield), a simple graph database for managing dependencies. Using Roaring Bitmaps, it provides very fast O(1) queries on large datasets, and `bomfactory` made this process simple.

If you are looking for thousands of SBOMs for testing and research, https://github.com/bitbomdev/bom-silo repository contains a large collection of SBOMs that you can use which was created by using `bomfactory`.

## What is bomfactory?

`bomfactory` is a powerful CLI tool designed to automate the downloading of Software Bill of Materials (SBOMs) for multiple repositories. SBOMs are crucial for software testing and security analysis, and `bomfactory` simplifies the process of obtaining them. This project draws inspiration from the [criticality_score](https://github.com/ossf/criticality_score) project to target critical projects. 

## Features

- **Download Criticality Score CSV**: Easily download the Criticality Score CSV file.
- **Load CSV into SQLite**: Load CSV data into an SQLite database for efficient querying.
- **Query Repositories**: Perform complex queries on the SQLite database to find repositories based on various criteria.
- **Download SBOMs**: Automatically download SBOMs for repositories matching your query.


## Quickstart

### Installation

Clone the repository and install the dependencies:

```sh
git clone https://github.com/bitbomdev/bomfactory.git
cd bomfactory
make build
```

### Usage

#### 1. Download the CSV file

```sh
bomfactory download-csv --url https://www.googleapis.com/download/storage/v1/b/ossf-criticality-score/o/2024.07.05%2F143335%2Fall.csv?generation=1721362287412491&alt=media --output data.csv
```

#### 2. Load the CSV data into SQLite

```sh
bomfactory load --csv data.csv --db data.db --start 1 --end 0
```

#### 3. Query the SQLite data

```sh
bomfactory query --filter "repo_language:==:Go" --filter "repo_star_count:>:100" --db data.db
```

#### 4. Download SBOMs for repositories

```sh
bomfactory download-sbom --filter "repo_language:==:Go" --token my_github_token --dir sbom_files --db data.db
```


---

Feel free to open an issue if you have any questions or need further assistance!