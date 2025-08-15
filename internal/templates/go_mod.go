package templates

// GoModTemplate generates go.mod file with enhanced dependencies
const GoModTemplate = `module {{.Project.Module}}

go 1.21

require (
{{- if eq .Server.Type "gin" }}
	github.com/gin-gonic/gin v1.9.1
{{- else if eq .Server.Type "fiber" }}
	github.com/gofiber/fiber/v2 v2.52.0
{{- end }}
{{- if or (eq .Database.Type "postgres") (eq .Database.Type "both") }}
	github.com/lib/pq v1.10.9
{{- end }}
{{- if or (eq .Database.Type "mongodb") (eq .Database.Type "both") }}
	go.mongodb.org/mongo-driver v1.13.1
{{- end }}
{{- if .Database.Migrations }}
	github.com/golang-migrate/migrate/v4 v4.16.2
{{- end }}
{{- if .Events.Enabled }}
	github.com/streadway/amqp v1.1.0
{{- end }}
	github.com/joho/godotenv v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/go-playground/validator/v10 v10.16.0
)
`
