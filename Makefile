.PHONY: help build install test clean fmt lint features

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Ldflags for version injection
LDFLAGS := -X 'github.com/nicovandenhove/stardoc/internal/cli.Version=$(VERSION)' \
           -X 'github.com/nicovandenhove/stardoc/internal/cli.Commit=$(COMMIT)' \
           -X 'github.com/nicovandenhove/stardoc/internal/cli.Date=$(DATE)'

# Build flags
BUILD_FLAGS := -trimpath -ldflags "$(LDFLAGS) -s -w"

# Default target
help:
	@echo "Stardoc - Development Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build the binary"
	@echo "  make install     Install to GOPATH/bin"
	@echo "  make test        Run all tests"
	@echo "  make features    Run BDD feature tests"
	@echo "  make clean       Remove build artifacts"
	@echo "  make fmt         Format code"
	@echo "  make lint        Run linters"
	@echo ""

# Build the binary
build:
	@echo "Building stardoc $(VERSION)..."
	@go build $(BUILD_FLAGS) -o bin/stardoc ./cmd/stardoc
	@echo "✅ Binary created at bin/stardoc"

# Install to GOPATH/bin
install:
	@echo "Installing stardoc $(VERSION)..."
	@go install $(BUILD_FLAGS) ./cmd/stardoc
	@echo "✅ Installed to $(shell go env GOPATH)/bin/stardoc"

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run BDD feature tests with godog
features:
	@echo "Running BDD feature tests..."
	@go test -v ./features/...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf dist/
	@go clean
	@echo "✅ Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

# Run linters
lint:
	@echo "Running linters..."
	@golangci-lint run ./...
	@echo "✅ Lint complete"

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Development dependencies
deps:
	@echo "Installing development dependencies..."
	@go install github.com/cucumber/godog/cmd/godog@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Dependencies installed"

# Run the tool locally (for testing)
run:
	@go run ./cmd/stardoc $(ARGS)

# Create release builds for multiple platforms
release:
	@echo "Building releases..."
	@mkdir -p dist
	@GOOS=darwin GOARCH=amd64 go build -o dist/stardoc-darwin-amd64 ./cmd/stardoc
	@GOOS=darwin GOARCH=arm64 go build -o dist/stardoc-darwin-arm64 ./cmd/stardoc
	@GOOS=linux GOARCH=amd64 go build -o dist/stardoc-linux-amd64 ./cmd/stardoc
	@GOOS=linux GOARCH=arm64 go build -o dist/stardoc-linux-arm64 ./cmd/stardoc
	@GOOS=windows GOARCH=amd64 go build -o dist/stardoc-windows-amd64.exe ./cmd/stardoc
	@echo "✅ Release builds created in dist/"

# Initialize module dependencies
mod:
	@echo "Downloading module dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Module dependencies updated"
