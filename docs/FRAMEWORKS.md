# Supported Web Frameworks

The Clean Architecture Generator supports multiple web frameworks. You can choose your preferred framework by setting the `server.type` in your `cta.json` configuration file.

## Supported Frameworks

### Gin (github.com/gin-gonic/gin)

- **Type**: `"gin"`
- **Description**: Fast HTTP web framework with a martini-like API
- **Best for**: High-performance APIs, microservices
- **Features**: Middleware support, JSON validation, error management

**Example configuration:**
\`\`\`json
{
"server": {
"type": "gin",
"port": 8080
}
}
\`\`\`

### Fiber (github.com/gofiber/fiber/v2)

- **Type**: `"fiber"`
- **Description**: Express-inspired web framework built on top of Fasthttp
- **Best for**: High-performance applications, Express.js developers
- **Features**: Zero memory allocation router, built-in middleware, WebSocket support

**Example configuration:**
\`\`\`json
{
"server": {
"type": "fiber",
"port": 3000
}
}
\`\`\`

## Framework Differences

### Handler Signatures

**Gin Handlers:**
\`\`\`go
func (h *UserHandler) Create(c *gin.Context) {
c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}
\`\`\`

**Fiber Handlers:**
\`\`\`go
func (h *UserHandler) Create(c *fiber.Ctx) error {
return c.Status(fiber.StatusCreated).JSON(fiber.Map{
"message": "User created",
})
}
\`\`\`

### Route Registration

**Gin:**
\`\`\`go
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
userGroup := router.Group("/api/v1/users")
{
userGroup.POST("", h.Create)
userGroup.GET("/:id", h.GetByID)
}
}
\`\`\`

**Fiber:**
\`\`\`go
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
userGroup := app.Group("/api/v1/users")

    userGroup.Post("", h.Create)
    userGroup.Get("/:id", h.GetByID)

}
\`\`\`

### Middleware

**Gin:**
\`\`\`go
router.Use(gin.Logger())
router.Use(gin.Recovery())
\`\`\`

**Fiber:**
\`\`\`go
app.Use(logger.New())
app.Use(recover.New())
\`\`\`

## Performance Comparison

| Framework | Requests/sec | Memory Usage | Latency  |
| --------- | ------------ | ------------ | -------- |
| Gin       | ~47,000      | Medium       | Low      |
| Fiber     | ~100,000+    | Low          | Very Low |

_Note: Performance may vary based on application complexity and hardware._

## Choosing a Framework

### Choose Gin if:

- You prefer a mature, stable framework
- You need extensive middleware ecosystem
- You're building traditional REST APIs
- You want comprehensive documentation and community support

### Choose Fiber if:

- You need maximum performance
- You're familiar with Express.js patterns
- You want built-in features like WebSocket support
- You're building high-throughput applications

## Migration Between Frameworks

The generated code structure remains the same regardless of the framework choice. Only the HTTP layer (handlers) changes, making it easy to switch frameworks if needed.

To migrate:

1. Update `server.type` in your `cta.json`
2. Regenerate the project
3. Update any custom middleware or framework-specific code
