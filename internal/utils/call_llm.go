package utils

import (
	"fmt"
	"lcma/internal/ai"
	"lcma/internal/config"
	"os"
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
	prompt, err := BuildPromptWithContext(config.PromptTemplatePath, config.OutputPath)
	if err != nil {
		return err
	}
	response, err := CallLLM(prompt)
	if err != nil {
		return err
	}

	// Save response to file
	err = os.WriteFile(config.ReportPath, []byte(response), 0644)
	if err != nil {
		return fmt.Errorf("failed to write report to file: %w", err)
	}

	return nil
}
