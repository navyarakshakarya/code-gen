package templates

// MainTemplate generates main.go file with enhanced initialization
const MainTemplate = `package main

import (
	"fmt"
	"log"

{{- if eq .Server.Type "gin" }}
	"github.com/gin-gonic/gin"
{{- else if eq .Server.Type "fiber" }}
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
{{- end }}
{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
	"{{.Project.Module}}/internal/infrastructure/database/postgres"
{{- end }}
{{- if or (eq .Database.Type "mongodb") (eq .Database.Type "both") }}
	"{{.Project.Module}}/internal/infrastructure/database/mongodb"
{{- end }}
{{- if .Events.Enabled }}
	"{{.Project.Module}}/internal/infrastructure/events"
{{- end }}
	"{{.Project.Module}}/internal/usecase"
	"{{.Project.Module}}/internal/infrastructure/config"
	pkglogger "{{.Project.Module}}/pkg/logger"
	httphandler "{{.Project.Module}}/internal/handler/http"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	pkglogger.Init(cfg.LogLevel)

{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
	// Initialize PostgreSQL database
	pgDB, err := postgres.NewConnection(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgDB.Close()
{{- end }}

{{- if or (eq .Database.Type "mongodb") (eq .Database.Type "both") }}
	// Initialize MongoDB database
	mongoDB, err := mongodb.NewConnection(cfg.MongoURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Disconnect()
{{- end }}

{{- if .Events.Enabled }}
	// Initialize event bus
	eventBus, err := events.NewRabbitMQBus(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer eventBus.Close()
{{- end }}

	// Initialize repositories, use cases, and handlers
	// Example dependency injection:
	// Repositories
	{{- range .Domains }}
	{{- $domainName := .Name }}
	{{- range .Repositories }}
	{{- if eq .Database "postgres" }}
	{{.Name | lower}} := postgres.New{{.Name }}(pgDB)
	{{- else if eq .Database "mongodb" }}
	{{.Name | lower}} := mongodb.New{{.Name }}(mongoDB.Database)
	{{- end }}
	{{- end }}
	{{- end }}


	// Usecases
	{{- range .Domains }}
	{{- $domainName := .Name }}
	{{- range .UseCases }}
	{{. | lower}}UseCase := usecase.New{{.}}({{$domainName | lower}}Repository, pkglogger.GetLogger())
	{{- end }}
	{{- end }}


	// Handlers
	{{- range .Domains }}
	{{- range .Handlers }}
	{{.Name | lower}} := httphandler.New{{.Name}}({{range .UseCases}}{{. | lower }}UseCase, {{end}}pkglogger.GetLogger())
	{{- end }}
	{{- end }}

	pkglogger.Info("{{.Project.Name}} server starting on port {{.Server.Port}}")

{{- if eq .Server.Type "gin" }}
	// Setup Gin server
	server := gin.Default()

	// Register routes here
	{{- range .Domains }}
	{{- range .Handlers }}
	{{.Name | lower}}.RegisterRoutes(server)
	{{- end }}
	{{- end }}

	server.Run(fmt.Sprintf(":%d", cfg.Port))
{{- else if eq .Server.Type "fiber" }}
	// Setup Fiber server
	app := fiber.New(fiber.Config{
		AppName: "{{.Project.Name}}",
	})

	// Add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Register routes here
	{{- range .Domains }}
	{{- range .Handlers }}
	{{.Name |lower}}.RegisterRoutes(app)
	{{- end }}
	{{- end }}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
{{- end }}
}
`
