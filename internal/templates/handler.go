package templates

// Enhanced HandlerTemplate with automatic use case imports and injection
const HandlerTemplate = `package http

import (
{{- if eq .Config.Server.Type "gin" }}
	"github.com/gin-gonic/gin"
{{- else if eq .Config.Server.Type "fiber" }}
	"github.com/gofiber/fiber/v2"
	"{{$.Config.Project.Module}}/internal/usecase"
	"{{.Config.Project.Module}}/pkg/logger"
)
{{- end }}

// {{.Handler.Name}} handles HTTP requests for {{.Domain.Name}}
type {{.Handler.Name}} struct {
{{- range .Handler.UseCases }}
	{{. | lower}}UseCase *usecase.{{.}}
{{- end }}
	logger logger.Logger
}

// New{{.Handler.Name}} creates a new {{.Handler.Name}}
func New{{.Handler.Name}}(
{{- range .Handler.UseCases }}
	{{. | lower}}UseCase *usecase.{{.}},
{{- end }}
	logger logger.Logger,
) *{{.Handler.Name}} {
	return &{{.Handler.Name}}{
{{- range .Handler.UseCases }}
		{{. | lower}}UseCase: {{. | lower}}UseCase,
{{- end }}
		logger: logger,
	}
}

{{- if eq .Config.Server.Type "gin" }}
// RegisterRoutes registers routes for {{.Handler.Name}}
func (h *{{.Handler.Name}}) RegisterRoutes(router *gin.Engine) {
	{{.Domain.Name | lower}}Group := router.Group("/api/v1/{{.Domain.Name | lower}}s")
	{
{{- range .Handler.UseCases }}
{{- if contains . "Create" }}
		{{$.Domain.Name | lower}}Group.POST("", h.{{.}})
{{- else if contains . "Get" }}
		{{$.Domain.Name | lower}}Group.GET("/:id", h.{{.}})
{{- else if contains . "Update" }}
		{{$.Domain.Name | lower}}Group.PUT("/:id", h.{{.}})
{{- else if contains . "Delete" }}
		{{$.Domain.Name | lower}}Group.DELETE("/:id", h.{{.}})
{{- else if contains . "List" }}
		{{$.Domain.Name | lower}}Group.GET("", h.{{.}})
{{- end }}
{{- end }}
	}
}

{{- range .Handler.UseCases }}
// {{.}} handles the {{.}} use case
func (h *{{$.Handler.Name}}) {{.}}(c *gin.Context) {
{{- if contains . "GetByID" }}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	req := &usecase.{{.}}Request{ID: id}
{{- else }}
	var req usecase.{{.}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
{{- end }}

	resp, err := h.{{. | lower}}UseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("{{.}} failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

{{- if contains . "Create" }}
	c.JSON(http.StatusCreated, resp)
{{- else if contains . "Delete" }}
	c.JSON(http.StatusNoContent, nil)
{{- else }}
	c.JSON(http.StatusOK, resp)
{{- end }}
}
{{- end }}

{{- else if eq .Config.Server.Type "fiber" }}
// RegisterRoutes registers routes for {{.Handler.Name}}
func (h *{{.Handler.Name}}) RegisterRoutes(app *fiber.App) {
	{{.Domain.Name | lower}}Group := app.Group("/api/v1/{{.Domain.Name | lower}}s")
	
{{- range .Handler.UseCases }}
{{- if contains . "Create" }}
	{{$.Domain.Name | lower}}Group.Post("", h.{{.}})
{{- else if contains . "Get" }}
	{{$.Domain.Name | lower}}Group.Get("/:id", h.{{.}})
{{- else if contains . "Update" }}
	{{$.Domain.Name | lower}}Group.Put("/:id", h.{{.}})
{{- else if contains . "Delete" }}
	{{$.Domain.Name | lower}}Group.Delete("/:id", h.{{.}})
{{- else if contains . "List" }}
	{{$.Domain.Name | lower}}Group.Get("", h.{{.}})
{{- end }}
{{- end }}
}

{{- range .Handler.UseCases }}
// {{.}} handles the {{.}} use case
func (h *{{$.Handler.Name}}) {{.}}(c *fiber.Ctx) error {
{{- if contains . "GetByID" }}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	req := &usecase.{{.}}Request{ID: id}
{{- else }}
	var req usecase.{{.}}Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
{{- end }}

	resp, err := h.{{. | lower}}UseCase.Execute(c.Context(), &req)
	if err != nil {
		h.logger.Error("{{.}} failed:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

{{- if contains . "Create" }}
	return c.Status(fiber.StatusCreated).JSON(resp)
{{- else if contains . "Delete" }}
	return c.SendStatus(fiber.StatusNoContent)
{{- else }}
	return c.JSON(resp)
{{- end }}
}
{{- end }}
{{- end }}
`
