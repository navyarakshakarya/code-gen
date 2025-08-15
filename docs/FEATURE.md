# Enhanced Features

## Database Support

The generator now supports multiple database options:

### PostgreSQL

\`\`\`json
{
"database": {
"type": "postgres",
"migrations": true,
"host": "localhost",
"port": 5432,
"name": "myapp"
}
}
\`\`\`

### MongoDB

\`\`\`json
{
"database": {
"type": "mongodb",
"migrations": false,
"host": "localhost",
"port": 27017,
"name": "myapp"
}
}
\`\`\`

### Both Databases

\`\`\`json
{
"database": {
"type": "both",
"migrations": true,
"host": "localhost",
"port": 5432,
"name": "myapp"
}
}
\`\`\`

## Event-Driven Architecture

Enable event-driven architecture with RabbitMQ:

\`\`\`json
{
"events": {
"enabled": true,
"type": "rabbitmq",
"host": "localhost",
"port": 5672
}
}
\`\`\`

## Enhanced Handler Configuration

Handlers now automatically import and wire up their required use cases:

\`\`\`json
{
"handlers": [
{
"name": "UserHandler",
"usecases": ["CreateUser", "GetUser", "UpdateUser", "DeleteUser"]
},
{
"name": "ProductHandler",
"usecases": ["CreateProduct", "GetProduct", "ListProducts"]
}
]
}
\`\`\`

This generates:

- Automatic use case imports
- Constructor with dependency injection
- Handler methods for each use case
- Proper framework-specific signatures (Gin/Fiber)

## Generated Infrastructure

The generator creates:

- `internal/infrastructure/database.go` - Database connections
- `internal/infrastructure/events.go` - Event bus implementation
- Proper dependency injection in handlers
- Framework-specific handler signatures
