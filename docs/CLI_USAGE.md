# CLI Usage Guide

## Commands

### `init` - Interactive Project Setup

Initialize a new clean architecture project with interactive prompts:

\`\`\`bash
go run main.go init
\`\`\`

This command will guide you through:

- **Project Information**: Name, Go module, description
- **Server Framework**: Choose between Gin or Fiber
- **Database Selection**: PostgreSQL, MongoDB, or both
- **Event System**: Enable/disable RabbitMQ event bus
- **Domain Configuration**: Define entities, repositories, use cases, and handlers

### `generate` or `gen` - Code Generation

Generate code from a configuration file:

\`\`\`bash

# Generate from default cta.json

go run main.go generate

# Generate from custom config file

go run main.go generate my-config.json

# Generate to specific output directory

go run main.go gen my-config.json ./output
\`\`\`

## Configuration Structure

The enhanced configuration supports:

### Repository Configuration

\`\`\`json
{
"repositories": [
{
"name": "UserRepository",
"database": "postgres",
"domain": "user"
},
{
"name": "ProfileRepository",
"database": "mongodb",
"domain": "user"
}
]
}
\`\`\`

### Handler Configuration

\`\`\`json
{
"handlers": [
{
"name": "UserHandler",
"usecases": ["CreateUser", "GetUser", "UpdateUser", "DeleteUser"]
}
]
}
\`\`\`

### Database Configuration

\`\`\`json
{
"database": {
"type": "both", // "postgres", "mongodb", or "both"
"host": "localhost",
"port": 5432,
"name": "myapp_db"
}
}
\`\`\`

### Event Configuration

\`\`\`json
{
"events": {
"enabled": true,
"type": "rabbitmq",
"url": "amqp://localhost:5672"
}
}
\`\`\`

## Generated Structure

\`\`\`
project/
├── cmd/server/ # Application entry point
├── internal/
│ ├── domain/
│ │ ├── entity/ # Business entities
│ │ └── repository/ # Repository interfaces
│ ├── usecase/ # Business logic
│ ├── handler/http/ # HTTP handlers
│ └── infrastructure/
│ ├── database/
│ │ ├── postgres/ # PostgreSQL implementations
│ │ └── mongodb/ # MongoDB implementations
│ ├── events/ # Event bus (if enabled)
│ └── config/ # Configuration
├── pkg/ # Shared packages
├── migrations/ # Database migrations
└── docs/ # Documentation
\`\`\`

## Features

- ✅ **Interactive Setup**: Guided project initialization
- ✅ **Multi-Database**: PostgreSQL, MongoDB, or both
- ✅ **Event-Driven**: RabbitMQ integration
- ✅ **Clean Architecture**: Proper layer separation
- ✅ **Domain-Driven**: Domain-specific organization
- ✅ **Auto-Wiring**: Automatic dependency injection
- ✅ **Framework Choice**: Gin or Fiber support
