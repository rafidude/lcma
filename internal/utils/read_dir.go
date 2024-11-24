package utils

import (
	"fmt"
	"lcma/internal/config"
	"os"
	"path/filepath"
	"strings"
)

// ReadDirectoryFiles reads all .py and .html files from the given directory
// and its subdirectories, combining their contents into a single output file
func ReadLegacyCodeGenerateOutput(dirPath string) error {
	// If dirPath is empty, read from env
	if dirPath == "" {
		dirPath = config.LegacyCodePath
	}

	// Create/truncate output file
	outputFile, err := os.Create(config.OutputFilePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	// Walk through directory
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		// Skip directories, including .venv directories
		if info.IsDir() {
			if filepath.Base(path) == ".venv" {
				return filepath.SkipDir // Skip this directory and all its contents
			}
			return nil
		}

		// Check if file extension is .py or .html
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".py" && ext != ".html" {
			return nil
		}

		// Read file contents
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Write filename header and contents to output file
		header := fmt.Sprintf("# %s\n", filepath.Base(path))
		if _, err := outputFile.WriteString(header); err != nil {
			return fmt.Errorf("error writing header to output file: %w", err)
		}

		if _, err := outputFile.Write(content); err != nil {
			return fmt.Errorf("error writing content to output file: %w", err)
		}

		// Add newline between files
		if _, err := outputFile.WriteString("\n\n"); err != nil {
			return fmt.Errorf("error writing newline to output file: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	return nil
}
