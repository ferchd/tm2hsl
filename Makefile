+.PHONY: all build test clean lint fmt vet mod-tidy install dev-setup release help

# Variables
BINARY_NAME=tm2hsl
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Default target
all: clean mod-tidy fmt vet lint test build

# Build the binary
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/tm2hsl

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/tm2hsl
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 ./cmd/tm2hsl
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/tm2hsl
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/tm2hsl
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/tm2hsl

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	go clean
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Tidy modules
mod-tidy:
	go mod tidy
	go mod verify

# Install dependencies
install-deps:
	go mod download

# Install the binary
install: build
	go install $(LDFLAGS) ./cmd/tm2hsl

# Development setup
dev-setup: install-deps
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

# Run all checks (CI)
ci: mod-tidy fmt vet lint test build

# Create release archives
release: build-all
	mkdir -p dist
	cp bin/* dist/
	cd dist && for f in *; do sha256sum "$$f" > "$$f.sha256"; done
	cd dist && tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd dist && tar -czf $(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	cd dist && tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd dist && tar -czf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	cd dist && zip $(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe

# Show version
version:
	@echo $(VERSION)

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Run clean, mod-tidy, fmt, vet, lint, test, and build"
	@echo "  build        - Build the binary"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Lint code"
	@echo "  mod-tidy     - Tidy modules"
	@echo "  install      - Install the binary"
	@echo "  dev-setup    - Set up development environment"
	@echo "  ci           - Run all CI checks"
	@echo "  release      - Create release archives"
	@echo "  version      - Show version"
	@echo "  help         - Show this help"