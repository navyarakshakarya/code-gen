package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/navyarakshakarya/code-gen/internal/config"
	"github.com/navyarakshakarya/code-gen/internal/templates"
)

// Execute executes a template with the given data
func Execute(tmpl string, data interface{}) (string, error) {
	t, err := template.New("template").Funcs(templates.GetTemplateFuncs()).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Generate creates the clean architecture boilerplate
func Generate(configPath, outputDir string) error {
	// Read and parse configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Create project structure
	if err := createProjectStructure(outputDir, cfg); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Generate files
	generator := &Generator{
		config:    cfg,
		outputDir: outputDir,
		options: GeneratorOptions{
			SkipExisting: true, // Default to skip existing files
			CreateBackup: false,
			Force:        false,
		},
	}

	return generator.generateAll()
}

// GenerateWithOptions creates the clean architecture boilerplate with options
func GenerateWithOptions(configPath, outputDir string, options GeneratorOptions) error {
	// Read and parse configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Create project structure
	if err := createProjectStructure(outputDir, cfg); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Generate files with options
	generator := &Generator{
		config:    cfg,
		outputDir: outputDir,
		options:   options,
	}

	return generator.generateAll()
}

// GeneratorOptions handles options for file generation
type GeneratorOptions struct {
	SkipExisting bool // Skip files that already exist
	CreateBackup bool // Create backup of existing files before overwriting
	Force        bool // Force overwrite existing files
}

// Generator handles code generation
type Generator struct {
	config    config.Config
	outputDir string
	options   GeneratorOptions
}

// generateAll generates all required files
func (g *Generator) generateAll() error {
	// Generate main files
	if err := g.generateMainFiles(); err != nil {
		return err
	}

	// Generate domain files
	for _, domain := range g.config.Domains {
		if err := g.generateDomainFiles(domain); err != nil {
			return err
		}
	}

	// Generate infrastructure files
	if err := g.generateInfrastructureFiles(); err != nil {
		return err
	}

	return nil
}

// generateMainFiles generates main application files
func (g *Generator) generateMainFiles() error {
	// Generate go.mod
	if err := g.writeFile("go.mod", templates.GoModTemplate, g.config); err != nil {
		return err
	}

	// Generate main.go
	if err := g.writeFile("cmd/server/main.go", templates.MainTemplate, g.config); err != nil {
		return err
	}

	// Generate README.md
	if err := g.writeFile("README.md", templates.ReadmeTemplate, g.config); err != nil {
		return err
	}

	// Generate Makefile
	if err := g.writeFile("Makefile", templates.MakefileTemplate, g.config); err != nil {
		return err
	}

	return nil
}

// generateDomainFiles generates files for each domain with enhanced handler support
func (g *Generator) generateDomainFiles(domain config.Domain) error {
	// Generate entities based on repository database types
	entityDatabaseMap := g.getEntityDatabaseMap(domain)

	for _, entity := range domain.Entities {
		// Determine database type for this entity
		dbType := entityDatabaseMap[entity]
		if dbType == "" {
			dbType = "postgres" // default to postgres
		}

		data := struct {
			Config       config.Config
			Domain       config.Domain
			Entity       string
			DatabaseType string
		}{g.config, domain, entity, dbType}

		filename := fmt.Sprintf("internal/domain/entity/%s.go", templates.ToSnakeCase(entity))

		// Use appropriate template based on database type
		var entityTemplate string
		if dbType == "mongodb" {
			entityTemplate = templates.MongoEntityTemplate
		} else {
			entityTemplate = templates.PostgresEntityTemplate
		}

		if err := g.writeFile(filename, entityTemplate, data); err != nil {
			return err
		}
	}

	if g.hasPostgresEntities(domain) {
		if err := g.generateSqlcFiles(domain); err != nil {
			return err
		}
	}

	for _, repo := range domain.Repositories {
		// Generate repository interface
		data := struct {
			Config     config.Config
			Domain     config.Domain
			Repository config.Repository
		}{g.config, domain, repo}

		filename := fmt.Sprintf("internal/domain/repository/%s.go", templates.ToSnakeCase(repo.Name))
		if err := g.writeFile(filename, templates.RepositoryTemplate, data); err != nil {
			return err
		}

		// Generate repository implementation based on database type
		implData := struct {
			Config     config.Config
			Domain     config.Domain
			Repository config.Repository
		}{g.config, domain, repo}

		var implTemplate string
		switch repo.Database {
		case "postgres":
			implTemplate = templates.PostgresRepositoryTemplate
		case "mongodb":
			implTemplate = templates.MongoRepositoryTemplate
		default:
			implTemplate = templates.PostgresRepositoryTemplate // default to postgres
		}

		implFilename := fmt.Sprintf("internal/infrastructure/database/%s/%s_impl.go",
			repo.Database, templates.ToSnakeCase(repo.Name))
		if err := g.writeFile(implFilename, implTemplate, implData); err != nil {
			return err
		}
	}

	// Generate use cases
	for _, usecase := range domain.UseCases {
		data := struct {
			Config  config.Config
			Domain  config.Domain
			UseCase string
		}{g.config, domain, usecase}

		filename := fmt.Sprintf("internal/usecase/%s.go", templates.ToSnakeCase(usecase))
		if err := g.writeFile(filename, templates.UseCaseTemplate, data); err != nil {
			return err
		}
	}

	for _, handler := range domain.Handlers {
		data := struct {
			Config  config.Config
			Domain  config.Domain
			Handler config.Handler
		}{g.config, domain, handler}

		filename := fmt.Sprintf("internal/handler/http/%s.go", templates.ToSnakeCase(handler.Name))
		if err := g.writeFile(filename, templates.HandlerTemplate, data); err != nil {
			return err
		}
	}

	return nil
}

// generateInfrastructureFiles generates infrastructure files with enhanced database and event support
func (g *Generator) generateInfrastructureFiles() error {
	// Generate database connections based on type
	if g.config.Database.Type == "postgres" || g.config.Database.Type == "both" {
		if err := g.writeFile("internal/infrastructure/database/postgres/connection.go", templates.PostgresTemplate, g.config); err != nil {
			return err
		}
		if err := g.writeFile("internal/infrastructure/database/postgres/with_schema.go", templates.WithConnSchema, g.config); err != nil {
			return err
		}
	}

	if g.config.Database.Type == "mongodb" || g.config.Database.Type == "both" {
		if err := g.writeFile("internal/infrastructure/database/mongodb/connection.go", templates.MongoTemplate, g.config); err != nil {
			return err
		}
	}

	// Generate legacy database connection for backward compatibility
	if g.config.Database.Type == "" {
		if err := g.writeFile("internal/infrastructure/database/connection.go", templates.PostgresTemplate, g.config); err != nil {
			return err
		}
	}

	// Generate event bus if enabled
	if g.config.Events.Enabled {
		if err := g.writeFile("internal/infrastructure/events/eventbus.go", templates.EventBusTemplate, g.config); err != nil {
			return err
		}
	}

	// Generate enhanced config or legacy config
	configTemplate := templates.ConfigTemplate
	if g.hasEnhancedFeatures() {
		configTemplate = templates.ConfigTemplate
	}

	if err := g.writeFile("internal/infrastructure/config/config.go", configTemplate, g.config); err != nil {
		return err
	}

	// Generate logger
	if err := g.writeFile("pkg/logger/logger.go", templates.LoggerTemplate, g.config); err != nil {
		return err
	}

	return nil
}

// hasEnhancedFeatures checks if the config uses enhanced features
func (g *Generator) hasEnhancedFeatures() bool {
	return g.config.Events.Enabled ||
		g.config.Database.Type == "mongodb" ||
		g.config.Database.Type == "both" ||
		(g.config.Database.Type == "postgres" && (g.config.Database.Host != "" || g.config.Database.Port != 0))
}

// writeFile writes content to a file using template with options
func (g *Generator) writeFile(filename, template string, data interface{}) error {
	fullPath := filepath.Join(g.outputDir, filename)

	// Check if file already exists
	if _, err := os.Stat(fullPath); err == nil {
		// File exists
		if g.options.SkipExisting && !g.options.Force {
			fmt.Printf("Skipping existing file: %s\n", filename)
			return nil
		}

		// Create backup if requested
		if g.options.CreateBackup {
			backupPath := fullPath + ".backup"
			if err := g.createBackup(fullPath, backupPath); err != nil {
				fmt.Printf("Warning: Failed to create backup for %s: %v\n", filename, err)
			} else {
				fmt.Printf("Created backup: %s.backup\n", filename)
			}
		}

		if !g.options.Force {
			fmt.Printf("Overwriting existing file: %s\n", filename)
		}
	}

	content, err := Execute(template, data)
	if err != nil {
		return fmt.Errorf("failed to execute template for %s: %w", filename, err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return nil
}

// createBackup creates a backup of an existing file
func (g *Generator) createBackup(srcPath, backupPath string) error {
	srcData, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath, srcData, 0644)
}

// isFileModified checks if a file has been modified from its original content
func (g *Generator) isFileModified(filename, template string, data interface{}) (bool, error) {
	fullPath := filepath.Join(g.outputDir, filename)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false, nil // File doesn't exist, so not modified
	}

	// Read existing file content
	existingContent, err := os.ReadFile(fullPath)
	if err != nil {
		return false, err
	}

	// Generate what the content should be
	expectedContent, err := Execute(template, data)
	if err != nil {
		return false, err
	}

	// Compare content (ignoring whitespace differences)
	existing := strings.TrimSpace(string(existingContent))
	expected := strings.TrimSpace(expectedContent)

	return existing != expected, nil
}

// createProjectStructure creates the folder structure with enhanced directories
func createProjectStructure(outputDir string, cfg config.Config) error {
	dirs := []string{
		"cmd/server",
		"internal/domain/entity",
		"internal/domain/repository",
		"internal/usecase",
		"internal/handler/http",
		"internal/infrastructure/config",
		"pkg/logger",
		"pkg/validator",
		"docs",
	}

	// Add database-specific directories
	if cfg.Database.Type == "postgres" || cfg.Database.Type == "both" {
		dirs = append(dirs,
			"internal/infrastructure/database/postgres",
			"internal/infrastructure/database/postgres/sqlc",
			"db/migrations",
			"db/queries",
		)
	}
	if cfg.Database.Type == "mongodb" || cfg.Database.Type == "both" {
		dirs = append(dirs, "internal/infrastructure/database/mongodb")
	}
	if cfg.Database.Type != "" {
		dirs = append(dirs, "internal/infrastructure/database")
	}

	// Add migrations directory if enabled
	if cfg.Database.Migrations {
		dirs = append(dirs, "migrations")
	}

	// Add events directory if enabled
	if cfg.Events.Enabled {
		dirs = append(dirs, "internal/infrastructure/events")
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(outputDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
	}

	return nil
}

// getEntityDatabaseMap maps entities to their respective database types
func (g *Generator) getEntityDatabaseMap(domain config.Domain) map[string]string {
	entityDBMap := make(map[string]string)

	for _, repo := range domain.Repositories {
		if repo.Entity != "" {
			entityDBMap[repo.Entity] = repo.Database
		}
	}

	return entityDBMap
}

// hasPostgresEntities checks if a domain has any PostgreSQL entities
func (g *Generator) hasPostgresEntities(domain config.Domain) bool {
	entityDBMap := g.getEntityDatabaseMap(domain)

	for _, dbType := range entityDBMap {
		if dbType == "postgres" {
			return true
		}
	}

	return false
}

// generateSqlcFiles generates SQLC configuration and queries for PostgreSQL entities
func (g *Generator) generateSqlcFiles(domain config.Domain) error {
	// Generate sqlc.yaml configuration
	if err := g.writeFile("sqlc.yaml", templates.SqlcConfigTemplate, g.config); err != nil {
		return err
	}

	// Generate queries for each PostgreSQL entity
	entityDBMap := g.getEntityDatabaseMap(domain)

	for _, entity := range domain.Entities {
		if entityDBMap[entity] == "postgres" {
			data := struct {
				Config config.Config
				Domain config.Domain
				Entity string
			}{g.config, domain, entity}

			filename := fmt.Sprintf("db/queries/%s.sql", templates.ToSnakeCase(entity))
			if err := g.writeFile(filename, templates.SqlcQueryTemplate, data); err != nil {
				return err
			}
		}
	}

	return nil
}
