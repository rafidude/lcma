package utils

import (
	"fmt"
	"lcma/internal/ai"
	"lcma/internal/config"
	"os"
	"path/filepath"
)

func CallLLM(prompt string) (string, error) {
	client := ai.NewGroqClient(config.GroqAPIKey)
	model := config.Model

	messages := []ai.GroqMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := client.CreateChatCompletion(model, messages)
	if err != nil {
		return "", err
	}

	// Print the response
	return response.Choices[0].Message.Content, nil
}

func CallLLMWithContextAndSaveReport() error {
	// Define file pairs for processing
	filePairs := []struct {
		promptFile string
		reportFile string
	}{
		{
			promptFile: "prompt.txt",
			reportFile: "report.md",
		},
		{
			promptFile: "prompt_code.txt",
			reportFile: "report_code.md",
		},
		{
			promptFile: "prompt_ui.txt",
			reportFile: "report_ui.md",
		},
	}

	outputFile, err := os.ReadFile(config.OutputFilePath)
	if err != nil {
		return fmt.Errorf("failed to read output file: %w", err)
	}

	// Process each file pair
	for _, pair := range filePairs {
		fmt.Println("Processing file pair:", pair)
		promptPath := filepath.Join(config.PromptTemplatePath, pair.promptFile)

		// Build prompt using the current pair of files
		prompt, err := buildPromptWithContext(promptPath)
		if err != nil {
			return fmt.Errorf("failed to build prompt for %s: %w", promptPath, err)
		}

		prompt += "\n\n<legacy_code>\n" + string(outputFile) + "\n</legacy_code>"

		// Call LLM with the constructed prompt
		response, err := CallLLM(prompt)
		if err != nil {
			return fmt.Errorf("failed to get LLM response for %s: %w", promptPath, err)
		}

		reportPath := filepath.Join(config.ReportPath, pair.reportFile)
		// Create report directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(reportPath), 0755); err != nil {
			return fmt.Errorf("failed to create report directory for %s: %w", pair.reportFile, err)
		}

		// Save response to corresponding report file
		if err := os.WriteFile(reportPath, []byte(response), 0644); err != nil {
			return fmt.Errorf("failed to write report %s: %w", pair.reportFile, err)
		}
	}

	return nil
}
