package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/navyarakshakarya/code-gen/analyzer"
	"github.com/navyarakshakarya/code-gen/generator"
	"github.com/navyarakshakarya/code-gen/logger"
)

const (
	version = "v1.0.0"
	banner  = `
 ██████╗ ██████╗ ██████╗ ███████╗      ██████╗ ███████╗███╗   ██╗
██╔════╝██╔═══██╗██╔══██╗██╔════╝     ██╔════╝ ██╔════╝████╗  ██║
██║     ██║   ██║██║  ██║█████╗       ██║  ███╗█████╗  ██╔██╗ ██║
██║     ██║   ██║██║  ██║██╔══╝       ██║   ██║██╔══╝  ██║╚██╗██║
╚██████╗╚██████╔╝██████╔╝███████╗     ╚██████╔╝███████╗██║ ╚████║
 ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝      ╚═════╝ ╚══════╝╚═╝  ╚═══╝

Go Clean Architecture Code Generator %s
`
)

func main() {
	var (
		verbose   = flag.Bool("verbose", false, "enable verbose output")
		version   = flag.Bool("version", false, "show version")
		help      = flag.Bool("help", false, "show help")
		dryRun    = flag.Bool("dry-run", false, "show what would be generated")
		force     = flag.Bool("force", false, "overwrite existing .gen.go files")
		tags      = flag.String("tags", "", "build tags to include")
		outputDir = flag.String("output", "", "output directory (default: current directory)")
	)

	flag.Parse()

	if *version {
		fmt.Printf("code-gen %v\n", version)
		return
	}

	if *help {
		printUsage()
		return
	}

	// Initialize logger
	logger := logger.New(*verbose)

	if *verbose {
		fmt.Printf("%v %v", banner, version)
	}

	// Get current working directory
	workDir, err := os.Getwd()
	if err != nil {
		logger.Fatal("Failed to get current directory: %v", err)
	}

	// Validate Go project
	if err := validateGoProject(workDir); err != nil {
		logger.Fatal("Invalid Go project: %v", err)
	}

	logger.Info("Analyzing Go project in: %s", workDir)

	// Initialize analyzer with build tags
	analyzer := analyzer.New(logger, *tags)

	// Analyze project
	projectInfo, err := analyzer.AnalyzeProject(workDir)
	if err != nil {
		logger.Fatal("Analysis failed: %v", err)
	}

	if len(projectInfo.Interfaces) == 0 {
		logger.Warning("No interfaces found in project")
		logger.Info("Make sure your interfaces follow naming conventions (e.g., *Repo, *UseCase, *Handler)")
		return
	}

	logger.Success("Analysis complete: found %d interfaces, %d structs",
		len(projectInfo.Interfaces), len(projectInfo.Structs))

	// Initialize generator
	gen := generator.New(logger)

	// Generate code
	results, err := gen.Generate(projectInfo)
	if err != nil {
		logger.Fatal("Code generation failed: %v", err)
	}

	// Determine output directory
	outDir := workDir
	if *outputDir != "" {
		outDir = *outputDir
	}

	// Write files or show dry run
	if *dryRun {
		logger.Info("Dry run - files that would be generated:")
		for _, result := range results {
			logger.Info("  %s (%d lines)", result.Filename, result.LineCount)
		}
	} else {
		written, skipped := writeFiles(results, outDir, *force, logger)

		logger.Success("Code generation complete!")
		logger.Info("Generated %d files, skipped %d existing files", written, skipped)

		if skipped > 0 {
			logger.Info("Use -force to overwrite existing files")
		}

		logger.Info("\nNext steps:")
		logger.Info("  1. Review generated code")
		logger.Info("  2. Implement TODO methods")
		logger.Info("  3. Run: go mod tidy")
		logger.Info("  4. Run: go build")
	}
}

func printUsage() {
	fmt.Printf(`code-gen - Go Clean Architecture Code Generator

USAGE:
    code-gen [flags]

FLAGS:
    -verbose        Enable verbose output
    -version        Show version information
    -help           Show this help message
    -dry-run        Show what would be generated without creating files
    -force          Overwrite existing .gen.go files
    -tags string    Build tags to include during analysis
    -output string  Output directory (default: current directory)

EXAMPLES:
    code-gen                    # Generate code for current project
    code-gen -verbose           # Enable verbose output
    code-gen -dry-run           # Preview what would be generated
    code-gen -force             # Overwrite existing files
    code-gen -tags "integration,dev"  # Include build tags

INSTALLATION:
    go install github.com/your-org/code-gen@latest

For more information, visit: https://github.com/your-org/code-gen
`)
}

func validateGoProject(dir string) error {
	// Check for go.mod
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		return fmt.Errorf("go.mod not found - not a Go module")
	}

	// Check for .go files
	hasGoFiles := false
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".go" && !info.IsDir() {
			hasGoFiles = true
			return filepath.SkipDir // Found at least one, can stop
		}
		return nil
	})

	if err != nil {
		return err
	}

	if !hasGoFiles {
		return fmt.Errorf("no Go source files found")
	}

	return nil
}

func writeFiles(results []*generator.GeneratedFile, outputDir string, force bool, logger *logger.Logger) (written, skipped int) {
	for _, result := range results {
		filePath := filepath.Join(outputDir, result.Filename)

		// Check if file exists
		if _, err := os.Stat(filePath); err == nil && !force {
			logger.Warning("File exists, skipping: %s", result.Filename)
			skipped++
			continue
		}

		// Create directory if needed
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			logger.Error("Failed to create directory: %v", err)
			continue
		}

		// Write file
		if err := os.WriteFile(filePath, []byte(result.Content), 0644); err != nil {
			logger.Error("Failed to write %s: %v", result.Filename, err)
			continue
		}

		logger.Success("Generated: %s", result.Filename)
		written++
	}

	return written, skipped
}
