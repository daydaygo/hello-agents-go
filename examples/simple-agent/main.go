package main

import (
	"context"
	"fmt"
	"os"

	"github.com/daydaygo/hello-agents-go/pkg/agent"
	"github.com/daydaygo/hello-agents-go/pkg/llm"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY environment variable")
		return
	}

	client := llm.NewClient(apiKey)
	ag := agent.NewSimpleAgent("assistant", client,
		agent.WithSystemPrompt("You are a helpful assistant. Be concise and friendly."),
	)

	fmt.Println("Simple Agent Example")
	fmt.Println("====================")
	fmt.Print("\nEnter your question: ")

	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	response, err := ag.Run(context.Background(), input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nResponse: %s\n", response)
}
