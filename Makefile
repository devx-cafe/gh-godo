# Makefile for gh-godo
.PHONY: help build test lint clean install coverage coverage-check watch build-all fmt vet

# Default target
.DEFAULT_GOAL := help

# Binary name
BINARY_NAME=godo

# Build directory
BUILD_DIR=.

# Coverage threshold (percentage)
COVERAGE_THRESHOLD=70

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Built $(BINARY_NAME)"

build-all: ## Build for multiple platforms (Linux, Darwin, Windows for AMD64 and ARM64)
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-arm64.exe .
	@echo "✅ Built binaries for all platforms in dist/"

test: ## Run tests
	$(GOTEST) -v ./...

coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.txt -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"
	@$(MAKE) coverage-check

coverage-check: ## Check if coverage meets the threshold
	@echo "Checking coverage against threshold ($(COVERAGE_THRESHOLD)%)..."
	@coverage=$$($(GOCMD) tool cover -func=coverage.txt | grep total | awk '{print $$3}' | sed 's/%//'); \
	threshold=$(COVERAGE_THRESHOLD); \
	if [ -z "$$coverage" ]; then \
		echo "❌ Could not determine coverage"; \
		exit 1; \
	fi; \
	coverage_int=$$(printf "%.0f" $$coverage); \
	if [ $$coverage_int -ge $$threshold ]; then \
		echo "✅ Coverage $$coverage% meets threshold $$threshold%"; \
	else \
		echo "❌ Coverage $$coverage% is below threshold $$threshold%"; \
		exit 1; \
	fi

lint: ## Run linter (requires golangci-lint)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Run: make install-lint"; \
	fi

install-lint: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.61.0
	@echo "✅ golangci-lint installed"

fmt: ## Format code
	$(GOFMT) ./...
	@echo "✅ Code formatted"

vet: ## Run go vet
	$(GOVET) ./...
	@echo "✅ Vet completed"

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -rf dist
	rm -f coverage.txt coverage.html
	@echo "✅ Cleaned build artifacts"

install: build ## Install the binary to $GOPATH/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/
	@echo "✅ Installed $(BINARY_NAME) to $(shell go env GOPATH)/bin"

watch: ## Watch for changes and rebuild (requires entr)
	@if command -v entr >/dev/null 2>&1; then \
		echo "👀 Watching for changes... (Press Ctrl+C to stop)"; \
		find . -name '*.go' | entr -c make build; \
	else \
		echo "⚠️  entr not installed. Install it with: apt-get install entr (Linux) or brew install entr (Mac)"; \
	fi

tidy: ## Tidy dependencies
	$(GOMOD) tidy
	@echo "✅ Dependencies tidied"

deps: ## Download dependencies
	$(GOMOD) download
	@echo "✅ Dependencies downloaded"

