package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"lcma/internal/config"
	"os"
	"path/filepath"
	"strings"
)

// CreateProjectStructure reads the report.md file and creates the modern project structure
func CreateProjectStructure(reportPath string) error {
	targetPath := config.ModernCodePath
	// Read the report file
	content, err := os.ReadFile(reportPath)
	if err != nil {
		return fmt.Errorf("failed to read report file: %w", err)
	}

	// Create the base project directory
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Extract and create directory structure
	if err := createDirectoryStructure(content, targetPath); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Extract and create Go files
	if err := createGoFiles(content, targetPath); err != nil {
		return fmt.Errorf("failed to create Go files: %w", err)
	}

	return nil
}

func createDirectoryStructure(content []byte, targetPath string) error {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	inJsonBlock := false
	var jsonContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "```json") {
			inJsonBlock = true
			continue
		}

		if line == "```" && inJsonBlock {
			// Parse the JSON and create directories
			var structure map[string]interface{}
			if err := json.Unmarshal([]byte(jsonContent.String()), &structure); err != nil {
				return fmt.Errorf("failed to parse JSON structure: %w", err)
			}

			// Create directories based on the JSON structure
			for _, dir := range []string{"models", "repositories", "services", "handlers", "templates", "public"} {
				fullPath := filepath.Join(targetPath, dir)
				if err := os.MkdirAll(fullPath, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
				}
			}
			break
		}

		if inJsonBlock {
			jsonContent.WriteString(line + "\n")
		}
	}
	return scanner.Err()
}

func createGoFiles(content []byte, targetPath string) error {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var currentFile string
	var fileContent strings.Builder
	inGoBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		// Look for Go code blocks
		if strings.HasPrefix(line, "**") && strings.HasSuffix(line, ".go**") {
			currentFile = strings.TrimSuffix(strings.TrimPrefix(line, "**"), "**")
			continue
		}

		if strings.HasPrefix(line, "```go") {
			inGoBlock = true
			fileContent.Reset()
			continue
		}

		if line == "```" && inGoBlock {
			inGoBlock = false
			if currentFile != "" {
				fullPath := filepath.Join(targetPath, currentFile)
				dir := filepath.Dir(fullPath)

				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
				}

				if err := os.WriteFile(fullPath, []byte(fileContent.String()), 0644); err != nil {
					return fmt.Errorf("failed to write file %s: %w", fullPath, err)
				}
				currentFile = ""
			}
			continue
		}

		if inGoBlock {
			fileContent.WriteString(line + "\n")
		}
	}
	return scanner.Err()
}
