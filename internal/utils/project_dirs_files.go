package utils

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/fs"
// 	"lcma/internal/config"
// 	"os"
// 	"path/filepath"
// 	"regexp"
// 	"strings"
// )

// // ProjectStructure represents the JSON structure of the project
// type ProjectStructure struct {
// 	Project struct {
// 		Name      string                 `json:"name"`
// 		Structure map[string]interface{} `json:"structure"`
// 	} `json:"project"`
// }

// // CodeBlock represents a parsed code block from the input
// type CodeBlock struct {
// 	Filename string
// 	Content  string
// }

// func CreateProjectStructure(reportPath string, targetPath string) error {
// 	// Read input file
// 	content, err := os.ReadFile(reportPath)
// 	if err != nil {
// 		return fmt.Errorf("error reading input file: %w", err)
// 	}

// 	// Parse JSON structure
// 	structure, err := parseJSONStructure(content)
// 	if err != nil {
// 		return fmt.Errorf("error parsing JSON structure: %w", err)
// 	}

// 	// Extract code blocks
// 	codeBlocks := extractCodeBlocks(string(content))
// 	// Print the number of code blocks extracted
// 	fmt.Printf("Number of code blocks extracted: %d\n", len(codeBlocks))

// 	// Create project directory
// 	projectDir := config.ModernCodePath
// 	if err := os.MkdirAll(projectDir, 0755); err != nil {
// 		return fmt.Errorf("error creating project directory: %w", err)
// 	}

// 	// Create directory structure and files
// 	if err := createProjectStructure(projectDir, structure.Project.Structure, codeBlocks); err != nil {
// 		return fmt.Errorf("error creating project structure: %w", err)
// 	}

// 	return nil
// }

// func parseJSONStructure(content []byte) (*ProjectStructure, error) {
// 	// Find and extract the JSON content from a ```json code block using regex
// 	re := regexp.MustCompile(`(?s)` + "```" + `json\n(.*?)` + "```")
// 	jsonMatch := re.FindSubmatch(content)
// 	if jsonMatch == nil {
// 		return nil, fmt.Errorf("no JSON content found in a ```json code block")
// 	}

// 	// Parse the JSON structure
// 	var structure ProjectStructure
// 	if err := json.Unmarshal(jsonMatch[1], &structure); err != nil {
// 		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
// 	}

// 	return &structure, nil
// }

// func extractCodeBlocks(content string) map[string]string {
// 	blocks := make(map[string]string)

// 	// Regular expression to match code block headers and content
// 	re := regexp.MustCompile(`(?m)^\*\*(.*?)\*\*\n` + "```" + `(?:go|html)\n(.*?)\n` + "```")
// 	matches := re.FindAllStringSubmatch(content, -1)

// 	for _, match := range matches {
// 		if len(match) >= 3 {
// 			filename := strings.TrimSpace(match[1])
// 			content := strings.TrimSpace(match[2])
// 			blocks[filename] = content
// 		}
// 	}

// 	// Extract code blocks from inline code blocks
// 	inlineCodeBlocks := extractInlineCodeBlocks(content)
// 	for filename, content := range inlineCodeBlocks {
// 		blocks[filename] = content
// 	}

// 	return blocks
// }

// func createProjectStructure(baseDir string, structure map[string]interface{}, codeBlocks map[string]string) error {
// 	return filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		relativePath := strings.TrimPrefix(path, baseDir)
// 		if relativePath == "" {
// 			// Process root directory
// 			return createDirectoryContents(baseDir, structure, codeBlocks)
// 		}

// 		return nil
// 	})
// }

// func createDirectoryContents(currentPath string, structure map[string]interface{}, codeBlocks map[string]string) error {
// 	for name, content := range structure {
// 		path := filepath.Join(currentPath, name)

// 		switch v := content.(type) {
// 		case string:
// 			// This is a file
// 			if codeBlock, exists := codeBlocks[path]; exists {
// 				if err := os.WriteFile(path, []byte(codeBlock), 0644); err != nil {
// 					return fmt.Errorf("error writing file %s: %v", path, err)
// 				}
// 				fmt.Printf("Created file: %s\n", path)
// 			}

// 		case map[string]interface{}:
// 			// This is a directory
// 			if err := os.MkdirAll(path, 0755); err != nil {
// 				return fmt.Errorf("error creating directory %s: %v", path, err)
// 			}
// 			fmt.Printf("Created directory: %s\n", path)

// 			// Recursively process directory contents
// 			if err := createDirectoryContents(path, v, codeBlocks); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func extractInlineCodeBlocks(content string) map[string]string {
// 	blocks := make(map[string]string)

// 	// Match inline code blocks with format: `filename:content`
// 	re := regexp.MustCompile("```([^:\n]+):([^\n]+)\n(.*?)\n```")
// 	matches := re.FindAllStringSubmatch(content, -1)

// 	for _, match := range matches {
// 		if len(match) >= 4 {
// 			filename := strings.TrimSpace(match[1])
// 			content := strings.TrimSpace(match[3])
// 			blocks[filename] = content
// 		}
// 	}

// 	return blocks
// }
