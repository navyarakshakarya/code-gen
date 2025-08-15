package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/navyarakshakarya/code-gen/internal/config"
	"github.com/navyarakshakarya/code-gen/internal/generator"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		if err := initProject(); err != nil {
			log.Fatalf("Error initializing project: %v", err)
		}
	case "generate", "gen":
		generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
		skipExisting := generateCmd.Bool("skip-existing", true, "Skip files that already exist (default: true)")
		createBackup := generateCmd.Bool("backup", false, "Create backup of existing files before overwriting")
		force := generateCmd.Bool("force", false, "Force overwrite existing files without prompting")

		// Parse remaining arguments after the command
		generateCmd.Parse(os.Args[2:])

		configPath := "cta.json"
		outputDir := "."

		args := generateCmd.Args()
		if len(args) > 0 {
			configPath = args[0]
		}
		if len(args) > 1 {
			outputDir = args[1]
		}

		options := generator.GeneratorOptions{
			SkipExisting: *skipExisting,
			CreateBackup: *createBackup,
			Force:        *force,
		}

		if err := generator.GenerateWithOptions(configPath, outputDir, options); err != nil {
			log.Fatalf("Error generating code: %v", err)
		}
		fmt.Println("✅ Clean architecture boilerplate generated successfully!")
	case "help", "-h", "--help":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		showHelp()
	}
}

func initProject() error {
	fmt.Println("🚀 Clean Architecture Generator - Project Initialization")
	fmt.Println("========================================================")

	reader := bufio.NewReader(os.Stdin)

	// Project basic info
	fmt.Print("📝 Project name: ")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	fmt.Print("📦 Go module (e.g., github.com/user/project): ")
	module, _ := reader.ReadString('\n')
	module = strings.TrimSpace(module)

	// Server framework selection
	fmt.Println("\n🌐 Server Framework:")
	fmt.Println("1. Gin")
	fmt.Println("2. Fiber")
	fmt.Print("Choose framework (1-2): ")
	frameworkChoice, _ := reader.ReadString('\n')
	frameworkChoice = strings.TrimSpace(frameworkChoice)

	var serverType string
	switch frameworkChoice {
	case "1":
		serverType = "gin"
	case "2":
		serverType = "fiber"
	default:
		serverType = "gin"
		fmt.Println("⚠️  Invalid choice, defaulting to Gin")
	}

	// Database selection
	fmt.Println("\n🗄️  Database:")
	fmt.Println("1. PostgreSQL")
	fmt.Println("2. MongoDB")
	fmt.Println("3. Both PostgreSQL and MongoDB")
	fmt.Print("Choose database (1-3): ")
	dbChoice, _ := reader.ReadString('\n')
	dbChoice = strings.TrimSpace(dbChoice)

	var dbType string
	switch dbChoice {
	case "1":
		dbType = "postgres"
	case "2":
		dbType = "mongodb"
	case "3":
		dbType = "both"
	default:
		dbType = "postgres"
		fmt.Println("⚠️  Invalid choice, defaulting to PostgreSQL")
	}

	// Event system
	fmt.Print("\n📡 Enable event system (RabbitMQ)? (y/N): ")
	eventChoice, _ := reader.ReadString('\n')
	eventChoice = strings.TrimSpace(strings.ToLower(eventChoice))
	enableEvents := eventChoice == "y" || eventChoice == "yes"

	// How many domains do you want to create? (1-5)
	fmt.Print("\n🏗️  How many domains do you want to create? (1-5): ")
	domainCountStr, _ := reader.ReadString('\n')
	domainCountStr = strings.TrimSpace(domainCountStr)
	domainCount, err := strconv.Atoi(domainCountStr)
	if err != nil || domainCount < 1 || domainCount > 5 {
		domainCount = 2
		fmt.Println("⚠️  Invalid number, defaulting to 2 domains")
	}

	domains := generateDefaultDomains(domainCount, dbType)

	fmt.Printf("\n✨ Generated %d default domains with example entities and CRUD operations:\n", len(domains))
	for i, domain := range domains {
		fmt.Printf("   %d. %s Domain - Entities: %s\n", i+1, domain.Name, strings.Join(domain.Entities, ", "))
	}

	// Create configuration
	cfg := config.Config{
		Project: config.Project{
			Name:        projectName,
			Module:      module,
			Description: fmt.Sprintf("Clean architecture application with %d domains", len(domains)),
		},
		Server: config.Server{
			Type: serverType,
			Port: 8080,
		},
		Database: config.Database{
			Type: dbType,
			Host: "localhost",
			Port: 5432,
			Name: projectName + "_db",
		},
		Events: config.Events{
			Enabled: enableEvents,
			Type:    "rabbitmq",
			Host:    "localhost",
			Port:    5672,
		},
		Domains: domains,
	}

	// Save configuration
	configFile, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile("cta.json", configFile, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	fmt.Println("\n✅ Configuration saved to cta.json")
	fmt.Println("🚀 Run 'go run main.go generate' to generate your project!")
	fmt.Println("📝 You can edit cta.json to customize the generated domains before running generate")

	return nil
}

func generateDefaultDomains(count int, dbType string) []config.Domain {
	domainTemplates := []struct {
		name     string
		entities []string
		usecases []string
	}{
		{
			name:     "User",
			entities: []string{"User", "Profile"},
			usecases: []string{"CreateUser", "GetUser", "UpdateUser", "DeleteUser", "GetUserProfile", "UpdateProfile"},
		},
		{
			name:     "Product",
			entities: []string{"Product", "Category"},
			usecases: []string{"CreateProduct", "GetProduct", "UpdateProduct", "DeleteProduct", "ListProducts", "GetProductsByCategory"},
		},
		{
			name:     "Order",
			entities: []string{"Order", "OrderItem"},
			usecases: []string{"CreateOrder", "GetOrder", "UpdateOrder", "CancelOrder", "ListOrders", "GetOrderHistory"},
		},
		{
			name:     "Auth",
			entities: []string{"Session", "Token"},
			usecases: []string{"Login", "Logout", "RefreshToken", "ValidateToken", "ResetPassword"},
		},
		{
			name:     "Notification",
			entities: []string{"Notification", "Template"},
			usecases: []string{"SendNotification", "GetNotifications", "MarkAsRead", "CreateTemplate", "UpdateTemplate"},
		},
	}

	var domains []config.Domain

	for i := 0; i < count && i < len(domainTemplates); i++ {
		template := domainTemplates[i]

		var repositories []config.Repository
		for _, entity := range template.entities {
			repoDbType := dbType
			if dbType == "both" {
				// Alternate between postgres and mongodb for variety
				if i%2 == 0 {
					repoDbType = "postgres"
				} else {
					repoDbType = "mongodb"
				}
			}

			repositories = append(repositories, config.Repository{
				Name:     entity + "Repository",
				Entity:   entity, // Added explicit entity field
				Database: repoDbType,
			})
		}

		// Create handler with all use cases
		handlers := []config.Handler{
			{
				Name:     template.name + "Handler",
				UseCases: template.usecases,
			},
		}

		domains = append(domains, config.Domain{
			Name:         template.name,
			Entities:     template.entities,
			Repositories: repositories,
			UseCases:     template.usecases,
			Handlers:     handlers,
		})
	}

	return domains
}

func parseCommaSeparated(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}
	}

	parts := strings.Split(input, ",")
	var result []string
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func showHelp() {
	fmt.Println(`🚀 Clean Architecture Generator

Commands:
  init                 Initialize a new project with interactive setup
  generate [config]    Generate code from configuration file
  gen [config]         Alias for generate
  help                 Show this help message

Usage:
  go run main.go init                           # Interactive project setup
  go run main.go generate                       # Generate from cta.json (skip existing files)
  go run main.go generate my-config.json        # Generate from custom config
  go run main.go gen my-project.json ./output    # Generate to specific directory

File Protection Options:
  -skip-existing       Skip files that already exist (default: true)
  -backup              Create backup of existing files before overwriting
  -force               Force overwrite existing files without prompting
  -no-skip-existing    Disable skipping existing files

Examples:
  # Initialize new project
  go run main.go init

  # Generate code (skip existing files by default)
  go run main.go generate

  # Generate with backup of existing files
  go run main.go generate -backup

  # Force overwrite all files
  go run main.go generate -force -no-skip-existing

  # Generate with custom config and output directory
  go run main.go gen my-project.json ./output

Features:
  ✅ Clean Architecture patterns
  ✅ Multiple server frameworks (Gin, Fiber)
  ✅ Multiple databases (PostgreSQL, MongoDB, or both)
  ✅ Event-driven architecture (RabbitMQ)
  ✅ Domain-driven design
  ✅ Automatic dependency injection
  ✅ Repository pattern with database-specific implementations
  ✅ File protection to prevent overwriting existing code`)
}
