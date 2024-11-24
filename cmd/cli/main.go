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

	err = utils.ReadLegacyCodeGenerateOutput("")
	if err != nil {
		log.Fatal(err)
	}

	err = utils.CallLLMWithContextAndSaveReport()
	if err != nil {
		log.Fatal(err)
	}
	// reportFile := filepath.Join(config.ReportPath, "report_code.md")
	// err = utils.CreateProjectStructure(reportFile)
	// // err = utils.CreateProjectStructure(reportFile, config.ModernCodePath)
	// if err != nil {
	// 	log.Fatalf("Failed to create project structure: %v", err)
	// }
}
