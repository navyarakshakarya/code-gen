# Clean Architecture Guidelines

## Overview

This project follows Clean Architecture principles to ensure maintainability, testability, and separation of concerns.

## Layer Structure

### 1. Domain Layer (`internal/domain/`)

**Entities** (`internal/domain/entity/`)

- Core business objects
- Contains business rules and validation
- Independent of external concerns
- Example: User, Product, Order

**Repository Interfaces** (`internal/domain/repository/`)

- Defines contracts for data access
- Abstract interfaces, no implementation
- Used by use cases to access data

### 2. Use Case Layer (`internal/usecase/`)

- Contains application-specific business logic
- Orchestrates data flow between entities and repositories
- Independent of UI, database, and external services
- Each use case represents a single business operation

### 3. Interface Adapters Layer

**HTTP Handlers** (`internal/handler/http/`)

- Converts HTTP requests to use case calls
- Handles HTTP-specific concerns (routing, middleware)
- Converts use case responses to HTTP responses

**Repository Implementations** (`internal/infrastructure/database/`)

- Concrete implementations of repository interfaces
- Database-specific logic
- Implements data persistence

### 4. Infrastructure Layer (`internal/infrastructure/`)

- External concerns (database, web server, etc.)
- Configuration management
- Third-party integrations

## Dependency Rule

Dependencies must point inward:

- Infrastructure → Interface Adapters → Use Cases → Domain
- Inner layers should not depend on outer layers
- Use dependency injection to invert dependencies

## Adding New Features

### 1. Define the Entity

\`\`\`go
// internal/domain/entity/product.go
type Product struct {
ID int64
Name string
Description string
Price decimal.Decimal
CreatedAt time.Time
UpdatedAt time.Time
}

func (p \*Product) Validate() error {
if p.Name == "" {
return errors.New("product name is required")
}
return nil
}
\`\`\`

### 2. Create Repository Interface

\`\`\`go
// internal/domain/repository/product_repository.go
type ProductRepository interface {
Create(ctx context.Context, product *entity.Product) error
GetByID(ctx context.Context, id int64) (*entity.Product, error)
Update(ctx context.Context, product *entity.Product) error
Delete(ctx context.Context, id int64) error
List(ctx context.Context, limit, offset int) ([]*entity.Product, error)
}
\`\`\`

### 3. Implement Use Cases

\`\`\`go
// internal/usecase/create_product.go
type CreateProduct struct {
productRepo repository.ProductRepository
logger logger.Logger
}

func (uc *CreateProduct) Execute(ctx context.Context, req *CreateProductRequest) (\*CreateProductResponse, error) {
product := &entity.Product{
Name: req.Name,
Description: req.Description,
Price: req.Price,
}

    if err := product.Validate(); err != nil {
        return nil, err
    }

    if err := uc.productRepo.Create(ctx, product); err != nil {
        return nil, err
    }

    return &CreateProductResponse{Product: product}, nil

}
\`\`\`

### 4. Create HTTP Handler

\`\`\`go
// internal/handler/http/product_handler.go
func (h *ProductHandler) Create(c *gin.Context) {
var req CreateProductRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

    resp, err := h.createProductUC.Execute(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, resp)

}
\`\`\`

### 5. Wire Dependencies

\`\`\`go
// cmd/server/main.go
func main() {
// Infrastructure
db := database.NewConnection(cfg.DatabaseURL)

    // Repositories
    productRepo := postgres.NewProductRepository(db)

    // Use Cases
    createProductUC := usecase.NewCreateProduct(productRepo, logger)

    // Handlers
    productHandler := http.NewProductHandler(createProductUC, logger)

    // Routes
    productHandler.RegisterRoutes(router)

}
\`\`\`

## Best Practices

### 1. Dependency Injection

- Use constructor injection
- Inject interfaces, not concrete types
- Keep constructors simple

### 2. Error Handling

- Use custom error types for business errors
- Wrap errors with context
- Handle errors at appropriate layers

### 3. Testing

- Test each layer independently
- Use mocks for external dependencies
- Focus on business logic in use case tests

### 4. Validation

- Validate at entity level for business rules
- Validate at handler level for input format
- Use consistent validation patterns

### 5. Logging

- Log at appropriate levels
- Include context in log messages
- Use structured logging

## Configuration Management

Use environment variables for configuration:

\`\`\`go
type Config struct {
DatabaseURL string `env:"DATABASE_URL"`
Port int `env:"PORT" envDefault:"8080"`
LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}
\`\`\`

## Database Migrations

Store migrations in the `migrations/` directory:

\`\`\`sql
-- migrations/001_create_products_table.up.sql
CREATE TABLE products (
id SERIAL PRIMARY KEY,
name VARCHAR(255) NOT NULL,
description TEXT,
price DECIMAL(10,2) NOT NULL,
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP DEFAULT NOW()
);
\`\`\`

## Testing Strategy

### Unit Tests

- Test business logic in isolation
- Mock external dependencies
- Focus on edge cases and error conditions

### Integration Tests

- Test complete workflows
- Use test database
- Test actual HTTP endpoints

### Example Test Structure

\`\`\`go
func TestCreateProduct(t \*testing.T) {
// Arrange
mockRepo := &mocks.ProductRepository{}
uc := usecase.NewCreateProduct(mockRepo, logger)

    req := &CreateProductRequest{
        Name:  "Test Product",
        Price: decimal.NewFromFloat(99.99),
    }

    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    // Act
    resp, err := uc.Execute(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    mockRepo.AssertExpectations(t)

}
