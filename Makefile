.PHONY: build test test-coverage clean install lint

# Build the binary
build:
	go build -o go-covimerage ./cmd/go-covimerage

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f go-covimerage coverage.out coverage.html
	rm -f .coverage coverage.xml

# Install the binary
install:
	go install ./cmd/go-covimerage

# Run linters
lint:
	go vet ./...
	gofmt -s -w .
	go mod tidy
