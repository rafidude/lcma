package utils

import (
	"fmt"
	"lcma/internal/config"
	"os"
	"strings"
)

func BuildPromptWithContext(templatePath string, outputPath string) (string, error) {
	// Read the template file
	prompt, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	// Read the output file
	codeContext, err := os.ReadFile(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to read output file: %w", err)
	}

	// Convert prompt to string
	promptWithContext := string(prompt)

	// Replace all placeholder tags with actual prompt
	replacements := map[string]string{
		"<oldtech_stack></oldtech_stack>":       "<oldtech_stack>\n" + config.LegacyTechStack + "\n</oldtech_stack>",
		"<moderntech_stack></moderntech_stack>": "<moderntech_stack>\n" + config.ModernTechStack + "\n</moderntech_stack>",
		"<original_code></original_code>":       "<original_code>\n" + string(codeContext) + "\n</original_code>",
	}

	for placeholder, replacement := range replacements {
		promptWithContext = strings.Replace(promptWithContext, placeholder, replacement, 1)
	}

	return promptWithContext, nil
}
