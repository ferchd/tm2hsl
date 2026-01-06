#!/bin/bash
# Development setup script for tm2hsl

set -e

echo "Setting up tm2hsl development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "ERROR: Go version $GO_VERSION is too old. Please upgrade to Go $REQUIRED_VERSION or later."
    exit 1
fi

echo "OK: Go $GO_VERSION is installed"

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Install development tools
echo "Installing development tools..."

# golangci-lint for linting
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
else
    echo "OK: golangci-lint is already installed"
fi

# staticcheck for additional static analysis
if ! command -v staticcheck &> /dev/null; then
    echo "Installing staticcheck..."
    go install honnef.co/go/tools/cmd/staticcheck@latest
else
    echo "OK: staticcheck is already installed"
fi

# goreleaser for releases (optional)
if ! command -v goreleaser &> /dev/null; then
    echo "Installing goreleaser..."
    go install github.com/goreleaser/goreleaser@latest
else
    echo "OK: goreleaser is already installed"
fi

# Run initial checks
echo "Running initial checks..."
make ci

echo ""
echo "Development environment setup complete!"
echo ""
echo "Available commands:"
echo "  make build      - Build the binary"
echo "  make test       - Run tests"
echo "  make lint       - Run linter"
echo "  make fmt        - Format code"
echo "  make clean      - Clean build artifacts"
echo "  make help       - Show all available targets"
echo ""
echo "For more information, see docs/DEVELOPMENT.md"