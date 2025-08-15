package templates

// UseCaseTemplate remains unchanged
const UseCaseTemplate = `package usecase

import (
	"context"

	"{{.Config.Project.Module}}/internal/domain/repository"
	"{{.Config.Project.Module}}/pkg/logger"
)

// {{.UseCase}} handles {{.UseCase}} business logic
type {{.UseCase}} struct {
	{{.Domain.Name | lower}}Repo repository.{{.Domain.Name }}Repository
	logger       logger.Logger
}

// New{{.UseCase}} creates a new {{.UseCase}} use case
func New{{.UseCase}}({{.Domain.Name | lower}}Repo repository.{{.Domain.Name }}Repository, logger logger.Logger) *{{.UseCase}} {
	return &{{.UseCase}}{	
		{{.Domain.Name | lower}}Repo: {{.Domain.Name | lower}}Repo,
		logger:       logger,
	}
}

// Execute executes the {{.UseCase}} use case
func (uc *{{.UseCase}}) Execute(ctx context.Context, req *{{.UseCase}}Request) (*{{.UseCase}}Response, error) {
	// TODO: Implement {{.UseCase}} business logic
	uc.logger.Info("Executing {{.UseCase}}")
	
	return &{{.UseCase}}Response{}, nil
}

// {{.UseCase}}Request represents the request for {{.UseCase}}
type {{.UseCase}}Request struct {
	// TODO: Define request fields
}

// {{.UseCase}}Response represents the response for {{.UseCase}}
type {{.UseCase}}Response struct {
	// TODO: Define response fields
}
`
