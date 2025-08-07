.PHONY: build clean test install help

# Binary name
BINARY_NAME=go-nanda-sdk

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty)"

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) cmd/nanda-sdk/main.go
	@echo "Build complete!"

# Build for different platforms
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux cmd/nanda-sdk/main.go

build-macos:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-macos cmd/nanda-sdk/main.go

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME).exe cmd/nanda-sdk/main.go

build-all: build-linux build-macos build-windows
	@echo "Builds for all platforms complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux $(BINARY_NAME)-macos $(BINARY_NAME).exe nanda-sdk
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary for current platform"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-macos  - Build for macOS"
	@echo "  build-windows- Build for Windows"
	@echo "  build-all    - Build for all platforms"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  install      - Install binary to /usr/local/bin"
	@echo "  help         - Show this help message" 