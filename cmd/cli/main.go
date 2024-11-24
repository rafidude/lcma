package main

import (
	"log"

	"lcma/internal/config"
	"lcma/internal/utils"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	// utils.GenerateLegacyCodeContext()
	// utils.CallLLM("Who is the president of the United States?")
	// err = utils.CallLLMWithContextAndSaveReport()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = createProject()
	if err != nil {
		log.Fatal(err)
	}
}

func createProject() error {
	err := utils.CreateProjectStructure("./report.md")
	if err != nil {
		log.Fatalf("Failed to create project structure: %v", err)
	}
	return nil
}
