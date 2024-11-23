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
	err = utils.CallLLMWithContextAndSaveReport()
	if err != nil {
		log.Fatal(err)
	}

}
