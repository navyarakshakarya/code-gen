package templates

// Enhanced Config Template with database and event support
const ConfigTemplate = `package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
	PostgresURL string
{{- end }}
{{- if or (eq .Database.Type "mongodb") (eq .Database.Type "both") }}
	MongoURL    string
{{- end }}
{{- if .Events.Enabled }}
	RabbitMQURL string
{{- end }}
	Port        int
	LogLevel    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("PORT", "{{.Server.Port}}"))
	if err != nil {
		port = {{.Server.Port}}
	}

	config := &Config{
		Port:     port,
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
	config.PostgresURL = getEnv("POSTGRES_URL", "postgres://localhost/{{.Project.Name}}?sslmode=disable")
{{- end }}

{{- if or (eq .Database.Type "mongodb") (eq .Database.Type "both") }}
	config.MongoURL = getEnv("MONGO_URL", "mongodb://localhost:27017/{{.Project.Name}}")
{{- end }}

{{- if .Events.Enabled }}
	config.RabbitMQURL = getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
{{- end }}

	return config, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`
