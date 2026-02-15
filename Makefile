.PHONY: build test test-race run clean help

# Build the MarkerService binary
build:
	@echo "Building MarkerService..."
	@cd marker-service-go && go build -o ../bin/marker-service

# Run all tests
test:
	@echo "Running tests..."
	@cd marker-service-go && go test ./... -v

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@cd marker-service-go && go test -race ./...

# Run the service
run: build
	@echo "Starting MarkerService on :8080..."
	@./bin/marker-service

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@cd marker-service-go && go clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@cd marker-service-go && go mod download && go mod tidy

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@cd marker-service-go && swag init

# Display help
help:
	@echo "MarkerService Makefile targets:"
	@echo "  build       - Build the MarkerService binary"
	@echo "  test        - Run all tests"
	@echo "  test-race   - Run tests with race detector"
	@echo "  run         - Build and run the service"
	@echo "  clean       - Remove build artifacts"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  swagger     - Generate Swagger documentation"
	@echo "  help        - Display this help message"
