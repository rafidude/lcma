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

	// err = utils.ReadLegacyCodeGenerateOutput("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = utils.CallLLMWithContextAndSaveReport()
	if err != nil {
		log.Fatal(err)
	}

	// err = utils.CreateProjectStructure("./report.md")
	// if err != nil {
	// 	log.Fatalf("Failed to create project structure: %v", err)
	// }
}
