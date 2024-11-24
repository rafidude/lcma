package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	GroqAPIKey         string
	Model              string
	LegacyCodePath     string
	LegacyTechStack    string
	ModernTechStack    string
	PromptTemplatePath string
	OutputPath         string
	ReportPath         string
	ModernCodePath     string
)

// Init loads the environment variables and initializes the configuration
func Init() error {
	// Load .env file if it exists
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	// Load configuration from environment variables
	GroqAPIKey = os.Getenv("GROQ_API_KEY")
	if GroqAPIKey == "" {
		return fmt.Errorf("GROQ_API_KEY not set in .env file")
	}

	Model = os.Getenv("MODEL")
	if Model == "" {
		return fmt.Errorf("MODEL not set in .env file")
	}

	LegacyCodePath = os.Getenv("LEGACY_CODE_PATH")
	if LegacyCodePath == "" {
		return fmt.Errorf("LEGACY_CODE_PATH not set in .env file")
	}

	LegacyTechStack = os.Getenv("LEGACY_TECH_STACK")
	if LegacyTechStack == "" {
		return fmt.Errorf("LEGACY_TECH_STACK not set in .env file")
	}

	ModernTechStack = os.Getenv("MODERN_TECH_STACK")
	if ModernTechStack == "" {
		return fmt.Errorf("MODERN_TECH_STACK not set in .env file")
	}

	PromptTemplatePath = os.Getenv("PROMPT_TEMPLATE_PATH")
	if PromptTemplatePath == "" {
		return fmt.Errorf("PROMPT_TEMPLATE_PATH not set in .env file")
	}

	OutputPath = os.Getenv("OUTPUT_PATH")
	if OutputPath == "" {
		return fmt.Errorf("OUTPUT_PATH not set in .env file")
	}

	ReportPath = os.Getenv("REPORT_PATH")
	if ReportPath == "" {
		return fmt.Errorf("REPORT_PATH not set in .env file")
	}

	ModernCodePath = os.Getenv("MODERN_CODE_PATH")
	if ModernCodePath == "" {
		return fmt.Errorf("MODERN_CODE_PATH not set in .env file")
	}

	return nil
}
