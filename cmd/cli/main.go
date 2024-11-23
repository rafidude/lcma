package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"lcma/internal/ai"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	client := ai.NewGroqClient(os.Getenv("GROQ_API_KEY"))
	model := os.Getenv("MODEL")

	messages := []ai.GroqMessage{
		{
			Role:    "user",
			Content: "Explain the importance of fast language models",
		},
	}

	response, err := client.CreateChatCompletion(model, messages)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response
	fmt.Println(response.Choices[0].Message.Content)
}
