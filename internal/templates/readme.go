package templates

// ReadmeTemplate remains unchanged
const ReadmeTemplate = `# {{.Project.Name}}

{{.Project.Description}}

## Architecture

This project follows Clean Architecture principles with the following structure:

` + "```" + `
{{.Project.Name}}/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── domain/
│   │   ├── entity/       # Business entities
│   │   └── repository/   # Repository interfaces
│   ├── usecase/          # Business logic
│   ├── handler/
│   │   └── http/         # HTTP handlers
│   └── infrastructure/
│       ├── database/     # Database implementation
{{- if .Events.Enabled }}
│       ├── events/       # Event bus implementation
{{- end }}
│       └── config/       # Configuration
├── pkg/                  # Shared packages
├── migrations/           # Database migrations
└── docs/                 # Documentation
` + "```" + `

## Features

- **Clean Architecture**: Proper separation of concerns
- **Multiple Frameworks**: Support for Gin and Fiber
{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
- **PostgreSQL**: SQL database support
{{- end }}
{{- if or (eq .Database.Type "mongodb") (eq .Database.Type "both") }}
- **MongoDB**: NoSQL database support
{{- end }}
{{- if .Events.Enabled }}
- **Event-Driven**: RabbitMQ event bus integration
{{- end }}
- **Dependency Injection**: Constructor-based DI
- **Structured Logging**: JSON logging with levels

## Getting Started

### Prerequisites

- Go 1.21 or higher
{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
- PostgreSQL
{{- end }}
{{- if or (eq .Database.Type "mongodb") (eq .Database.Type "both") }}
- MongoDB
{{- end }}
{{- if .Events.Enabled }}
- RabbitMQ
{{- end }}

### Installation

1. Clone the repository:
` + "```bash" + `
git clone <repository-url>
cd {{.Project.Name}}
` + "```" + `

2. Install dependencies:
` + "```bash" + `
go mod tidy
` + "```" + `

3. Set up environment variables:
` + "```bash" + `
cp .env.example .env
# Edit .env with your configuration
` + "```" + `

4. Run the application:
` + "```bash" + `
make run
` + "```" + `

## Development

### Adding New Features

1. Define entities in ` + "`internal/domain/entity/`" + `
2. Create repository interfaces in ` + "`internal/domain/repository/`" + `
3. Implement use cases in ` + "`internal/usecase/`" + `
4. Create HTTP handlers in ` + "`internal/handler/http/`" + `
5. Wire dependencies in ` + "`cmd/server/main.go`" + `

### Running Tests

` + "```bash" + `
make test
` + "```" + `

### Building

` + "```bash" + `
make build
` + "```" + `

## API Documentation

API documentation is available at ` + "`/docs`" + ` when running the server.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request
`
