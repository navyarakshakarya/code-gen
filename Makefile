.PHONY: build clean test install

# Build the application
build:
	go build -o bin/code-gen main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Install the application
install:
	go install .

# Generate wire dependencies
generate:
	go generate ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run the application with example parameters
example:
	go run main.go -name example-project -template standard

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/go-boilerplate-generator-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/go-boilerplate-generator-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/go-boilerplate-generator-windows-amd64.exe main.go
