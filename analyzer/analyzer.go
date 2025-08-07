package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/navyaraskhakarya/code-gen/logger"
	"github.com/navyaraskhakarya/code-gen/types"
)

// Analyzer analyzes Go source code to extract interfaces and structs
type Analyzer struct {
	logger    *logger.Logger
	fileSet   *token.FileSet
	buildTags []string
}

// New creates a new analyzer instance
func New(logger *logger.Logger, tags string) *Analyzer {
	var buildTags []string
	if tags != "" {
		buildTags = strings.Split(tags, ",")
		for i, tag := range buildTags {
			buildTags[i] = strings.TrimSpace(tag)
		}
	}

	return &Analyzer{
		logger:    logger,
		fileSet:   token.NewFileSet(),
		buildTags: buildTags,
	}
}

// AnalyzeProject analyzes the entire Go project
func (a *Analyzer) AnalyzeProject(projectDir string) (*types.ProjectInfo, error) {
	projectInfo := &types.ProjectInfo{
		Interfaces: make(map[string]*types.InterfaceInfo),
		Structs:    make(map[string]*types.StructInfo),
		Imports:    make(map[string]string),
		ProjectDir: projectDir,
	}

	// Get module information
	moduleName, err := a.getModuleName(projectDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get module name: %w", err)
	}
	projectInfo.ModuleName = moduleName

	// Parse all Go files in the project
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and certain directories
		if !strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go") ||
			strings.HasSuffix(path, ".gen.go") ||
			strings.Contains(path, "vendor/") ||
			strings.Contains(path, ".git/") ||
			strings.Contains(path, "testdata/") {
			return nil
		}

		return a.analyzeFile(path, projectInfo)
	})

	if err != nil {
		return nil, err
	}

	// Post-process to establish relationships
	a.establishRelationships(projectInfo)

	return projectInfo, nil
}

// getModuleName extracts module name from go.mod
func (a *Analyzer) getModuleName(projectDir string) (string, error) {
	goModPath := filepath.Join(projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	return "", fmt.Errorf("module declaration not found in go.mod")
}

// analyzeFile analyzes a single Go file
func (a *Analyzer) analyzeFile(filePath string, projectInfo *types.ProjectInfo) error {
	// Check build constraints
	if !a.shouldIncludeFile(filePath) {
		a.logger.Info("Skipping file due to build constraints: %s", filePath)
		return nil
	}

	src, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	file, err := parser.ParseFile(a.fileSet, filePath, src, parser.ParseComments)
	if err != nil {
		a.logger.Warning("Failed to parse %s: %v", filePath, err)
		return nil // Continue with other files
	}

	packageName := file.Name.Name
	if projectInfo.PackageName == "" {
		projectInfo.PackageName = packageName
	}

	relPath, _ := filepath.Rel(projectInfo.ProjectDir, filePath)

	// Extract imports
	for _, imp := range file.Imports {
		if imp.Path != nil {
			importPath := strings.Trim(imp.Path.Value, `"`)
			var alias string
			if imp.Name != nil {
				alias = imp.Name.Name
			} else {
				// Extract package name from import path
				parts := strings.Split(importPath, "/")
				alias = parts[len(parts)-1]
			}
			projectInfo.Imports[alias] = importPath
		}
	}

	// Analyze declarations
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			if node.Tok == token.TYPE {
				a.processTypeDeclaration(node, packageName, relPath, projectInfo)
			}
		}
		return true
	})

	return nil
}

// shouldIncludeFile checks if file should be included based on build tags
func (a *Analyzer) shouldIncludeFile(filePath string) bool {
	if len(a.buildTags) == 0 {
		return true
	}

	// Read first few lines to check build constraints
	content, err := os.ReadFile(filePath)
	if err != nil {
		return true
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if i > 10 { // Only check first 10 lines
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//go:build") || strings.HasPrefix(line, "// +build") {
			// Simple build tag checking - in production, use go/build package
			for _, tag := range a.buildTags {
				if strings.Contains(line, tag) {
					return true
				}
			}
			return false
		}
	}

	return true
}

// processTypeDeclaration processes type declarations
func (a *Analyzer) processTypeDeclaration(genDecl *ast.GenDecl, packageName, filePath string, projectInfo *types.ProjectInfo) {
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		// Extract comments
		var comments []string
		if genDecl.Doc != nil {
			for _, comment := range genDecl.Doc.List {
				comments = append(comments, strings.TrimPrefix(comment.Text, "//"))
			}
		}

		switch t := typeSpec.Type.(type) {
		case *ast.InterfaceType:
			a.extractInterface(typeSpec.Name.Name, t, packageName, filePath, comments, projectInfo)
		case *ast.StructType:
			a.extractStruct(typeSpec.Name.Name, t, packageName, filePath, comments, projectInfo)
		}
	}
}

// extractInterface extracts interface information
func (a *Analyzer) extractInterface(name string, iface *ast.InterfaceType, pkg, filePath string, comments []string, projectInfo *types.ProjectInfo) {
	interfaceInfo := &types.InterfaceInfo{
		Name:     name,
		Package:  pkg,
		FilePath: filePath,
		Methods:  []types.MethodInfo{},
		Layer:    a.determineLayer(name),
		Comments: comments,
	}

	// Extract methods
	for _, method := range iface.Methods.List {
		if funcType, ok := method.Type.(*ast.FuncType); ok {
			for _, methodName := range method.Names {
				methodInfo := a.extractMethodInfo(methodName.Name, funcType)
				interfaceInfo.Methods = append(interfaceInfo.Methods, methodInfo)
			}
		}
	}

	projectInfo.Interfaces[name] = interfaceInfo
	a.logger.Info("Found interface: %s (%s layer)", name, interfaceInfo.Layer)
}

// extractStruct extracts struct information
func (a *Analyzer) extractStruct(name string, structType *ast.StructType, pkg, filePath string, comments []string, projectInfo *types.ProjectInfo) {
	structInfo := &types.StructInfo{
		Name:     name,
		Package:  pkg,
		FilePath: filePath,
		Fields:   []types.FieldInfo{},
		Comments: comments,
	}

	// Extract fields
	for _, field := range structType.Fields.List {
		fieldType := a.typeToString(field.Type)
		var tag string
		if field.Tag != nil {
			tag = field.Tag.Value
		}

		if len(field.Names) > 0 {
			// Named fields
			for _, fieldName := range field.Names {
				structInfo.Fields = append(structInfo.Fields, types.FieldInfo{
					Name: fieldName.Name,
					Type: fieldType,
					Tag:  tag,
				})
			}
		} else {
			// Embedded field
			structInfo.Fields = append(structInfo.Fields, types.FieldInfo{
				Name:     "",
				Type:     fieldType,
				Tag:      tag,
				Embedded: true,
			})
		}
	}

	projectInfo.Structs[name] = structInfo
}

// extractMethodInfo extracts method information from function type
func (a *Analyzer) extractMethodInfo(name string, funcType *ast.FuncType) types.MethodInfo {
	method := types.MethodInfo{
		Name:    name,
		Params:  []types.ParamInfo{},
		Returns: []types.ParamInfo{},
	}

	// Extract parameters
	if funcType.Params != nil {
		for _, param := range funcType.Params.List {
			paramType := a.typeToString(param.Type)

			// Check for context.Context
			if strings.Contains(paramType, "Context") {
				method.HasContext = true
			}

			if len(param.Names) > 0 {
				for _, paramName := range param.Names {
					method.Params = append(method.Params, types.ParamInfo{
						Name: paramName.Name,
						Type: paramType,
					})
				}
			} else {
				method.Params = append(method.Params, types.ParamInfo{
					Name: "",
					Type: paramType,
				})
			}
		}
	}

	// Extract return types
	if funcType.Results != nil {
		for _, result := range funcType.Results.List {
			resultType := a.typeToString(result.Type)

			// Check for error return
			if resultType == "error" {
				method.HasError = true
			}

			if len(result.Names) > 0 {
				for _, resultName := range result.Names {
					method.Returns = append(method.Returns, types.ParamInfo{
						Name: resultName.Name,
						Type: resultType,
					})
				}
			} else {
				method.Returns = append(method.Returns, types.ParamInfo{
					Name: "",
					Type: resultType,
				})
			}
		}
	}

	return method
}

// typeToString converts AST type to string representation
func (a *Analyzer) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return a.typeToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + a.typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + a.typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + a.typeToString(t.Key) + "]" + a.typeToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		return "chan " + a.typeToString(t.Value)
	case *ast.FuncType:
		return "func(...)"
	case *ast.Ellipsis:
		return "..." + a.typeToString(t.Elt)
	default:
		return "interface{}"
	}
}

// determineLayer determines the architectural layer based on interface name
func (a *Analyzer) determineLayer(interfaceName string) types.LayerType {
	name := strings.ToLower(interfaceName)

	if strings.Contains(name, "repo") || strings.Contains(name, "repository") {
		return types.RepositoryLayer
	} else if strings.Contains(name, "usecase") || strings.Contains(name, "use_case") {
		return types.UseCaseLayer
	} else if strings.Contains(name, "handler") || strings.Contains(name, "controller") {
		return types.HandlerLayer
	} else if strings.Contains(name, "service") {
		return types.ServiceLayer
	}

	return types.ServiceLayer
}

// establishRelationships finds relationships between interfaces
func (a *Analyzer) establishRelationships(projectInfo *types.ProjectInfo) {
	for _, interfaceInfo := range projectInfo.Interfaces {
		baseName := a.extractBaseName(interfaceInfo.Name)

		// Find related interfaces with same base name
		for _, otherInterface := range projectInfo.Interfaces {
			otherBaseName := a.extractBaseName(otherInterface.Name)
			if baseName == otherBaseName && interfaceInfo.Name != otherInterface.Name {
				interfaceInfo.RelatedInterfaces = append(interfaceInfo.RelatedInterfaces, otherInterface.Name)
			}
		}
	}
}

// extractBaseName extracts the base name from interface name
func (a *Analyzer) extractBaseName(interfaceName string) string {
	suffixes := []string{"Handler", "Controller", "UseCase", "Service", "Repo", "Repository"}

	for _, suffix := range suffixes {
		if strings.HasSuffix(interfaceName, suffix) {
			return strings.TrimSuffix(interfaceName, suffix)
		}
	}

	return interfaceName
}
