.PHONY: help build test fmt vet clean coverage bench demo cli

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build the CLI tool"
	@echo "  test      - Run tests"
	@echo "  fmt       - Format Go code"
	@echo "  vet       - Run go vet"
	@echo "  clean     - Clean build artifacts"
	@echo "  coverage  - Run tests with coverage"
	@echo "  bench     - Run benchmarks"
	@echo "  demo      - Build WebAssembly demo"
	@echo "  cli       - Build and install CLI tool"
	@echo "  check     - Run all checks (fmt, vet, test)"

# Build the CLI tool
build:
	go build -o xatena-cli ./cmd/xatena-cli

# Run tests
test:
	go test -v ./...

# Format Go code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Clean build artifacts
clean:
	rm -f xatena-cli
	rm -f coverage*
	rm -f demo/main.wasm
	rm -f demo/wasm_exec.js

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out -coverpkg=./... ./pkg/xatena ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Total coverage:"
	@go tool cover -func=coverage.out | grep "total:"

# Run benchmarks
bench:
	./bench/bench.sh

# Build WebAssembly demo
demo:
	cd demo && $(MAKE) build

# Build and install CLI tool
cli: build
	@echo "CLI tool built as ./xatena-cli"

# Run all checks
check: fmt vet test
	@echo "All checks passed!"

# Development workflow
dev: fmt vet test
	@echo "Development checks completed!"

# Install dependencies (if needed)
deps:
	go mod tidy
	go mod download

# Run specific test
test-pkg:
	go test -v ./pkg/xatena

test-internal:
	go test -v ./internal/...

# Show module information
mod-info:
	go mod tidy
	go list -m all