BINARY_NAME=gh-deployer
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

.PHONY: build test clean install lint fmt vet

# Build the application
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Lint the code
lint:
	golangci-lint run

# Format the code
fmt:
	go fmt ./...

# Vet the code
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Install the binary
install: build
	sudo cp ${BINARY_NAME} /usr/local/bin/

# Create systemd service file
systemd-service:
	@echo "Creating systemd service file..."
	@echo "Copy the example service file from the project documentation"
	@echo "See .github/copilot-instructions.md for systemd service configuration"

# Development build with race detection
dev-build:
	go build -race ${LDFLAGS} -o ${BINARY_NAME} .

# Run in development mode with file watching
dev-run:
	./$(BINARY_NAME) --config config.yaml --dry-run