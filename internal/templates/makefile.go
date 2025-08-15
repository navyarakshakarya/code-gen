package templates

// MakefileTemplate remains unchanged
const MakefileTemplate = `# {{.Project.Name}} Makefile

.PHONY: build run test clean deps migrate-up migrate-down

# Build the application
build:
	go build -o bin/{{.Project.Name}} cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

{{- if .Database.Migrations }}
# Run database migrations up
migrate-up:
{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
	migrate -path migrations -database "$(POSTGRES_URL)" up
{{- end }}

# Run database migrations down
migrate-down:
{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
	migrate -path migrations -database "$(POSTGRES_URL)" down
{{- end }}
{{- end }}

# Generate code (if using code generation tools)
generate:
	go generate ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run the application in development mode
dev:
	air

# Docker commands
docker-build:
	docker build -t {{.Project.Name}} .

docker-run:
	docker run -p {{.Server.Port}}:{{.Server.Port}} {{.Project.Name}}

{{- if .Events.Enabled }}
# Start RabbitMQ for development
rabbitmq-start:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
{{- end }}

{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
# Start PostgreSQL for development
postgres-start:
	docker run -d --name postgres -p 5432:5432 -e POSTGRES_DB={{.Project.Name}} -e POSTGRES_PASSWORD=password postgres:15
{{- end }}

{{- if or (eq .Database.Type "mongodb") (eq .Database.Type "both") }}
# Start MongoDB for development
mongo-start:
	docker run -d --name mongodb -p 27017:27017 mongo:7
{{- end }}
`
