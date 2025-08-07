package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/navyarakshakarya/code-gen/logger"
	"github.com/navyarakshakarya/code-gen/types"
)

// Generator generates clean architecture code
type Generator struct {
	logger *logger.Logger
}

// GeneratedFile represents a generated file
type GeneratedFile struct {
	Filename  string
	Content   string
	LineCount int
}

// New creates a new generator instance
func New(logger *logger.Logger) *Generator {
	return &Generator{
		logger: logger,
	}
}

// Generate generates all code files
func (g *Generator) Generate(projectInfo *types.ProjectInfo) ([]*GeneratedFile, error) {
	var results []*GeneratedFile

	// Generate implementations for each interface
	for interfaceName, interfaceInfo := range projectInfo.Interfaces {
		file, err := g.generateImplementation(interfaceName, interfaceInfo, projectInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to generate implementation for %s: %w", interfaceName, err)
		}
		results = append(results, file)
	}

	// Generate factory
	factoryFile, err := g.generateFactory(projectInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to generate factory: %w", err)
	}
	results = append(results, factoryFile)

	// Generate wire integration (similar to Google Wire)
	wireFile, err := g.generateWireIntegration(projectInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to generate wire integration: %w", err)
	}
	results = append(results, wireFile)

	return results, nil
}

// generateImplementation generates implementation for an interface
func (g *Generator) generateImplementation(interfaceName string, interfaceInfo *types.InterfaceInfo, projectInfo *types.ProjectInfo) (*GeneratedFile, error) {
	structName := g.generateStructName(interfaceName)
	fileName := g.generateFileName(interfaceName, interfaceInfo.Layer)

	var content strings.Builder

	// File header
	g.writeFileHeader(&content, projectInfo.PackageName)

	// Imports
	imports := g.generateImports(interfaceInfo, projectInfo)
	if len(imports) > 0 {
		content.WriteString("import (\n")
		for _, imp := range imports {
			content.WriteString(fmt.Sprintf("\t%s\n", imp))
		}
		content.WriteString(")\n\n")
	}

	// Struct definition
	g.writeStructDefinition(&content, structName, interfaceName, interfaceInfo, projectInfo)

	// Constructor
	g.writeConstructor(&content, structName, interfaceName, interfaceInfo, projectInfo)

	// Method implementations
	for _, method := range interfaceInfo.Methods {
		g.writeMethodImplementation(&content, structName, method, interfaceInfo.Layer)
	}

	// Interface compliance check
	content.WriteString(fmt.Sprintf("// Ensure %s implements %s\n", structName, interfaceName))
	content.WriteString(fmt.Sprintf("var _ %s = (*%s)(nil)\n", interfaceName, structName))

	return &GeneratedFile{
		Filename:  fileName,
		Content:   content.String(),
		LineCount: strings.Count(content.String(), "\n"),
	}, nil
}

// generateFactory generates the dependency injection factory
func (g *Generator) generateFactory(projectInfo *types.ProjectInfo) (*GeneratedFile, error) {
	var content strings.Builder

	g.writeFileHeader(&content, projectInfo.PackageName)

	// Imports
	content.WriteString("import (\n")
	content.WriteString("\t\"database/sql\"\n")
	content.WriteString("\t\"context\"\n")
	content.WriteString(")\n\n")

	// Factory struct
	content.WriteString("// Factory provides centralized dependency injection\n")
	content.WriteString("// This follows the factory pattern for clean architecture\n")
	content.WriteString("type Factory struct {\n")
	content.WriteString("\tdb     *sql.DB\n")
	content.WriteString("\tctx    context.Context\n")
	content.WriteString("\tconfig *Config // Add your config struct\n")
	content.WriteString("}\n\n")

	// Factory constructor
	content.WriteString("// NewFactory creates a new factory instance\n")
	content.WriteString("func NewFactory(db *sql.DB, ctx context.Context, config *Config) *Factory {\n")
	content.WriteString("\treturn &Factory{\n")
	content.WriteString("\t\tdb:     db,\n")
	content.WriteString("\t\tctx:    ctx,\n")
	content.WriteString("\t\tconfig: config,\n")
	content.WriteString("\t}\n")
	content.WriteString("}\n\n")

	// Generate factory methods for each interface
	for interfaceName, interfaceInfo := range projectInfo.Interfaces {
		g.writeFactoryMethod(&content, interfaceName, interfaceInfo, projectInfo)
	}

	return &GeneratedFile{
		Filename:  "factory.gen.go",
		Content:   content.String(),
		LineCount: strings.Count(content.String(), "\n"),
	}, nil
}

// generateWireIntegration generates Wire-compatible provider functions
func (g *Generator) generateWireIntegration(projectInfo *types.ProjectInfo) (*GeneratedFile, error) {
	var content strings.Builder

	g.writeFileHeader(&content, projectInfo.PackageName)

	// Wire build constraint
	content.WriteString("//go:build wireinject\n")
	content.WriteString("// +build wireinject\n\n")

	// Imports
	content.WriteString("import (\n")
	content.WriteString("\t\"database/sql\"\n")
	content.WriteString("\t\"context\"\n")
	content.WriteString("\t\"github.com/google/wire\"\n")
	content.WriteString(")\n\n")

	// Provider set
	content.WriteString("// ProviderSet is the Wire provider set for dependency injection\n")
	content.WriteString("var ProviderSet = wire.NewSet(\n")

	for interfaceName := range projectInfo.Interfaces {
		constructorName := "New" + interfaceName
		content.WriteString(fmt.Sprintf("\t%s,\n", constructorName))
	}

	content.WriteString("\tNewFactory,\n")
	content.WriteString(")\n\n")

	// Wire injector functions
	for interfaceName, interfaceInfo := range projectInfo.Interfaces {
		if interfaceInfo.Layer == types.HandlerLayer {
			g.writeWireInjector(&content, interfaceName, projectInfo)
		}
	}

	return &GeneratedFile{
		Filename:  "wire.gen.go",
		Content:   content.String(),
		LineCount: strings.Count(content.String(), "\n"),
	}, nil
}

// Helper methods for code generation

func (g *Generator) writeFileHeader(content *strings.Builder, packageName string) {
	content.WriteString("// Code generated by code-gen. DO NOT EDIT.\n")
	content.WriteString(fmt.Sprintf("// Generated at: %s\n\n", time.Now().Format(time.RFC3339)))
	content.WriteString(fmt.Sprintf("package %s\n\n", packageName))
}

func (g *Generator) generateStructName(interfaceName string) string {
	return strings.ToLower(string(interfaceName[0])) + interfaceName[1:]
}

func (g *Generator) generateFileName(interfaceName string, layer types.LayerType) string {
	baseName := g.extractBaseName(interfaceName)
	return fmt.Sprintf("%s_%s.gen.go", strings.ToLower(baseName), layer)
}

func (g *Generator) generateImports(interfaceInfo *types.InterfaceInfo, projectInfo *types.ProjectInfo) []string {
	imports := make(map[string]bool)

	// Standard library imports
	for _, method := range interfaceInfo.Methods {
		if method.HasContext {
			imports["\"context\""] = true
		}
		if method.HasError {
			// error is built-in, no import needed
		}

		// Check for common framework imports
		for _, param := range method.Params {
			g.addFrameworkImports(param.Type, imports)
		}
		for _, ret := range method.Returns {
			g.addFrameworkImports(ret.Type, imports)
		}
	}

	// Layer-specific imports
	switch interfaceInfo.Layer {
	case types.RepositoryLayer:
		imports["\"database/sql\""] = true
		imports["\"fmt\""] = true
	case types.UseCaseLayer:
		imports["\"fmt\""] = true
	case types.HandlerLayer:
		imports["\"encoding/json\""] = true
		imports["\"net/http\""] = true
	}

	// Convert to slice
	var result []string
	for imp := range imports {
		result = append(result, imp)
	}

	return result
}

func (g *Generator) addFrameworkImports(typeName string, imports map[string]bool) {
	if strings.Contains(typeName, "fiber.Ctx") {
		imports["\"github.com/gofiber/fiber/v2\""] = true
	}
	if strings.Contains(typeName, "gin.Context") {
		imports["\"github.com/gin-gonic/gin\""] = true
	}
	if strings.Contains(typeName, "echo.Context") {
		imports["\"github.com/labstack/echo/v4\""] = true
	}
	if strings.Contains(typeName, "http.ResponseWriter") || strings.Contains(typeName, "http.Request") {
		imports["\"net/http\""] = true
	}
}

func (g *Generator) writeStructDefinition(content *strings.Builder, structName, interfaceName string, interfaceInfo *types.InterfaceInfo, projectInfo *types.ProjectInfo) {
	// Comments
	if len(interfaceInfo.Comments) > 0 {
		for _, comment := range interfaceInfo.Comments {
			content.WriteString(fmt.Sprintf("// %s\n", strings.TrimSpace(comment)))
		}
	}

	content.WriteString(fmt.Sprintf("// %s implements %s interface\n", structName, interfaceName))
	content.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	// Dependencies
	dependencies := g.generateDependencies(interfaceName, interfaceInfo, projectInfo)
	for _, dep := range dependencies {
		content.WriteString(fmt.Sprintf("\t%s\n", dep))
	}

	content.WriteString("}\n\n")
}

func (g *Generator) writeConstructor(content *strings.Builder, structName, interfaceName string, interfaceInfo *types.InterfaceInfo, projectInfo *types.ProjectInfo) {
	dependencies := g.generateDependencies(interfaceName, interfaceInfo, projectInfo)

	content.WriteString(fmt.Sprintf("// New%s creates a new instance of %s\n", interfaceName, structName))
	content.WriteString(fmt.Sprintf("func New%s(", interfaceName))

	// Parameters
	var params []string
	var assignments []string

	for _, dep := range dependencies {
		parts := strings.Fields(dep)
		if len(parts) >= 2 {
			fieldName := parts[0]
			fieldType := strings.Join(parts[1:], " ")
			params = append(params, fmt.Sprintf("%s %s", fieldName, fieldType))
			assignments = append(assignments, fmt.Sprintf("\t\t%s: %s,", fieldName, fieldName))
		}
	}

	content.WriteString(strings.Join(params, ", "))
	content.WriteString(fmt.Sprintf(") %s {\n", interfaceName))
	content.WriteString(fmt.Sprintf("\treturn &%s{\n", structName))

	for _, assignment := range assignments {
		content.WriteString(assignment + "\n")
	}

	content.WriteString("\t}\n")
	content.WriteString("}\n\n")
}

func (g *Generator) writeMethodImplementation(content *strings.Builder, structName string, method types.MethodInfo, layer types.LayerType) {
	// Method signature
	content.WriteString(fmt.Sprintf("// %s implements the %s method\n", method.Name, method.Name))
	content.WriteString(fmt.Sprintf("func (impl *%s) %s(", structName, method.Name))

	// Parameters
	var params []string
	for _, param := range method.Params {
		if param.Name != "" {
			params = append(params, fmt.Sprintf("%s %s", param.Name, param.Type))
		} else {
			params = append(params, param.Type)
		}
	}
	content.WriteString(strings.Join(params, ", "))
	content.WriteString(")")

	// Return types
	if len(method.Returns) > 0 {
		content.WriteString(" (")
		var returns []string
		for _, ret := range method.Returns {
			if ret.Name != "" {
				returns = append(returns, fmt.Sprintf("%s %s", ret.Name, ret.Type))
			} else {
				returns = append(returns, ret.Type)
			}
		}
		content.WriteString(strings.Join(returns, ", "))
		content.WriteString(")")
	}

	content.WriteString(" {\n")

	// Method body with layer-specific templates
	g.writeMethodBody(content, method, layer)

	content.WriteString("}\n\n")
}

func (g *Generator) writeMethodBody(content *strings.Builder, method types.MethodInfo, layer types.LayerType) {
	content.WriteString(fmt.Sprintf("\t// TODO: Implement %s\n", method.Name))

	switch layer {
	case types.RepositoryLayer:
		content.WriteString("\t// Example database operation:\n")
		content.WriteString("\t// query := \"SELECT * FROM table WHERE condition = ?\"\n")
		content.WriteString("\t// rows, err := impl.db.QueryContext(ctx, query, param)\n")
		content.WriteString("\t// if err != nil {\n")
		content.WriteString("\t//     return result, fmt.Errorf(\"database query failed: %w\", err)\n")
		content.WriteString("\t// }\n")
		content.WriteString("\t// defer rows.Close()\n")
	case types.UseCaseLayer:
		content.WriteString("\t// Example business logic:\n")
		content.WriteString("\t// 1. Validate input parameters\n")
		content.WriteString("\t// 2. Call repository methods\n")
		content.WriteString("\t// 3. Apply business rules\n")
		content.WriteString("\t// 4. Return processed result\n")
	case types.HandlerLayer:
		content.WriteString("\t// Example HTTP handler:\n")
		content.WriteString("\t// 1. Parse request parameters\n")
		content.WriteString("\t// 2. Call use case methods\n")
		content.WriteString("\t// 3. Handle errors appropriately\n")
		content.WriteString("\t// 4. Return HTTP response\n")
	}

	// Generate return statement
	if len(method.Returns) > 0 {
		var returnValues []string
		for _, ret := range method.Returns {
			returnValues = append(returnValues, g.generateZeroValue(ret.Type))
		}
		content.WriteString(fmt.Sprintf("\treturn %s\n", strings.Join(returnValues, ", ")))
	}
}

func (g *Generator) writeFactoryMethod(content *strings.Builder, interfaceName string, interfaceInfo *types.InterfaceInfo, projectInfo *types.ProjectInfo) {
	baseName := g.extractBaseName(interfaceName)

	content.WriteString(fmt.Sprintf("// New%s creates a new %s instance with dependencies\n", interfaceName, interfaceName))
	content.WriteString(fmt.Sprintf("func (f *Factory) New%s() %s {\n", interfaceName, interfaceName))

	switch interfaceInfo.Layer {
	case types.RepositoryLayer:
		content.WriteString(fmt.Sprintf("\treturn New%s(f.db)\n", interfaceName))
	case types.UseCaseLayer:
		repoInterface := g.findRelatedInterface(baseName, types.RepositoryLayer, projectInfo)
		if repoInterface != "" {
			content.WriteString(fmt.Sprintf("\trepo := f.New%s()\n", repoInterface))
			content.WriteString(fmt.Sprintf("\treturn New%s(repo)\n", interfaceName))
		} else {
			content.WriteString("\t// TODO: Add repository dependency\n")
			content.WriteString(fmt.Sprintf("\treturn New%s(/* dependencies */)\n", interfaceName))
		}
	case types.HandlerLayer:
		useCaseInterface := g.findRelatedInterface(baseName, types.UseCaseLayer, projectInfo)
		if useCaseInterface != "" {
			content.WriteString(fmt.Sprintf("\tuseCase := f.New%s()\n", useCaseInterface))
			content.WriteString(fmt.Sprintf("\treturn New%s(useCase)\n", interfaceName))
		} else {
			content.WriteString("\t// TODO: Add use case dependency\n")
			content.WriteString(fmt.Sprintf("\treturn New%s(/* dependencies */)\n", interfaceName))
		}
	default:
		content.WriteString(fmt.Sprintf("\treturn New%s()\n", interfaceName))
	}

	content.WriteString("}\n\n")
}

func (g *Generator) writeWireInjector(content *strings.Builder, interfaceName string, projectInfo *types.ProjectInfo) {
	content.WriteString(fmt.Sprintf("// Initialize%s creates a fully wired %s instance\n", interfaceName, interfaceName))
	content.WriteString(fmt.Sprintf("func Initialize%s(db *sql.DB, ctx context.Context, config *Config) (%s, error) {\n", interfaceName, interfaceName))
	content.WriteString("\twire.Build(ProviderSet)\n")
	content.WriteString("\treturn nil, nil // Wire will generate the implementation\n")
	content.WriteString("}\n\n")
}

// Helper methods

func (g *Generator) generateDependencies(interfaceName string, interfaceInfo *types.InterfaceInfo, projectInfo *types.ProjectInfo) []string {
	var deps []string
	baseName := g.extractBaseName(interfaceName)

	switch interfaceInfo.Layer {
	case types.RepositoryLayer:
		deps = append(deps, "db *sql.DB")
	case types.UseCaseLayer:
		repoInterface := g.findRelatedInterface(baseName, types.RepositoryLayer, projectInfo)
		if repoInterface != "" {
			deps = append(deps, fmt.Sprintf("repo %s", repoInterface))
		}
	case types.HandlerLayer:
		useCaseInterface := g.findRelatedInterface(baseName, types.UseCaseLayer, projectInfo)
		if useCaseInterface != "" {
			deps = append(deps, fmt.Sprintf("useCase %s", useCaseInterface))
		}
	}

	return deps
}

func (g *Generator) findRelatedInterface(baseName string, layer types.LayerType, projectInfo *types.ProjectInfo) string {
	suffixes := map[types.LayerType][]string{
		types.RepositoryLayer: {"Repo", "Repository"},
		types.UseCaseLayer:    {"UseCase", "Service"},
		types.HandlerLayer:    {"Handler", "Controller"},
	}

	for _, suffix := range suffixes[layer] {
		interfaceName := baseName + suffix
		if _, exists := projectInfo.Interfaces[interfaceName]; exists {
			return interfaceName
		}
	}

	return ""
}

func (g *Generator) generateZeroValue(typeName string) string {
	switch {
	case typeName == "error":
		return "nil"
	case typeName == "string":
		return `""`
	case strings.HasPrefix(typeName, "int") || strings.HasPrefix(typeName, "uint") ||
		typeName == "float32" || typeName == "float64":
		return "0"
	case typeName == "bool":
		return "false"
	case strings.HasPrefix(typeName, "*") || strings.HasPrefix(typeName, "[]") ||
		strings.HasPrefix(typeName, "map[") || strings.Contains(typeName, "interface"):
		return "nil"
	default:
		return fmt.Sprintf("%s{}", typeName)
	}
}

func (g *Generator) extractBaseName(interfaceName string) string {
	suffixes := []string{"Handler", "Controller", "UseCase", "Service", "Repo", "Repository"}

	for _, suffix := range suffixes {
		if strings.HasSuffix(interfaceName, suffix) {
			return strings.TrimSuffix(interfaceName, suffix)
		}
	}

	return interfaceName
}
