package config

// Config represents the main configuration structure
type Config struct {
	Project  Project  `json:"project"`
	Domains  []Domain `json:"domains"`
	Database Database `json:"database"`
	Server   Server   `json:"server"`
	Events   Events   `json:"events,omitempty"` // Added events configuration
}

// Project contains project-level configuration
type Project struct {
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

// Domain represents a business domain
type Domain struct {
	Name         string       `json:"name"`
	Entities     []string     `json:"entities"`
	Repositories []Repository `json:"repositories"` // Changed from []string to []Repository for database-specific repos
	UseCases     []string     `json:"usecases"`
	Handlers     []Handler    `json:"handlers"` // Changed from []string to []Handler
}

// Handler represents a handler with its associated use cases
type Handler struct {
	Name     string   `json:"name"`
	UseCases []string `json:"usecases"` // Added use cases mapping
}

// Repository represents a repository with its database configuration
type Repository struct {
	Name     string `json:"name"`     // Repository name (e.g., "UserRepository")
	Entity   string `json:"entity"`   // Added entity field to specify which entity this repository handles
	Database string `json:"database"` // Database type: "postgres", "mongodb"
}

// Database configuration
type Database struct {
	Type       string `json:"type"` // postgres, mongodb, or both
	Migrations bool   `json:"migrations"`
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Name       string `json:"name,omitempty"`
}

// Events configuration for event-driven architecture
type Events struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"` // rabbitmq, kafka, etc.
	Host    string `json:"host,omitempty"`
	Port    int    `json:"port,omitempty"`
}

// Server configuration
type Server struct {
	Type string `json:"type"`
	Port int    `json:"port"`
}
