# X-Pay Go SDK Makefile

.PHONY: test build examples clean install deps lint fmt vet

# Default target
all: fmt vet test

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run ./...

# Vet code
vet:
	go vet ./...

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Run integration tests (requires valid API keys)
test-integration:
	go test -v -tags=integration ./...

# Build all examples
build-examples:
	@echo "Building examples..."
	go build -o bin/basic_payment ./examples/basic_payment.go
	go build -o bin/customers ./examples/customers.go
	go build -o bin/webhooks ./examples/webhooks.go

# Run basic payment example
run-basic:
	go run ./examples/basic_payment.go

# Run customer management example
run-customers:
	go run ./examples/customers.go

# Run webhook example
run-webhooks:
	go run ./examples/webhooks.go

# Run SDK test suite
test-sdk:
	go run ./test_sdk.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install the SDK locally
install:
	go mod download

# Generate documentation
docs:
	godoc -http=:6060

# Run all examples in sequence
run-all-examples: run-basic run-customers run-webhooks

# Benchmark tests
benchmark:
	go test -bench=. ./...

# Build for different platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/basic_payment ./examples/basic_payment.go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/basic_payment ./examples/basic_payment.go
	GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/basic_payment.exe ./examples/basic_payment.go

# Help
help:
	@echo "Available targets:"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  vet            - Run go vet"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-sdk       - Run SDK test suite"
	@echo "  build-examples - Build all examples"
	@echo "  run-basic      - Run basic payment example"
	@echo "  run-customers  - Run customer management example"
	@echo "  run-webhooks   - Run webhook example"
	@echo "  run-all-examples - Run all examples"
	@echo "  clean          - Clean build artifacts"
	@echo "  docs           - Generate documentation"
	@echo "  help           - Show this help"