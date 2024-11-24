package utils

import (
	"fmt"
	"lcma/internal/config"
	"os"
	"strings"
)

func buildPromptWithContext(templatePath string) (string, error) {
	// Read the template file
	prompt, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	// Convert prompt to string
	promptWithContext := string(prompt)

	// Replace all placeholder tags with actual prompt
	replacements := map[string]string{
		"<legacytech_stack></legacytech_stack>": "<legacytech_stack>\n" + config.LegacyTechStack + "\n</legacytech_stack>",
		"<moderntech_stack></moderntech_stack>": "<moderntech_stack>\n" + config.ModernTechStack + "\n</moderntech_stack>",
	}

	for placeholder, replacement := range replacements {
		promptWithContext = strings.Replace(promptWithContext, placeholder, replacement, 1)
	}

	return promptWithContext, nil
}
