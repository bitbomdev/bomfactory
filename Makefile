# Set the shell to bash
SHELL := /bin/bash

# Default target
all: build

# Build the project
build:
	go build -o bin/bom-factory main.go

# Lint the code
lint:
	golangci-lint run ./...

# Clean the build
clean:
	rm -rf bin/

.PHONY: all build lint clean
