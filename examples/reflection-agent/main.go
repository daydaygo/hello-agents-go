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
	ag := agent.NewReflectionAgent("reflector", client,
		agent.WithMaxReflections(2),
	)

	fmt.Println("Reflection Agent Example")
	fmt.Println("========================")
	fmt.Println("This agent will reflect on its work and improve it.")
	fmt.Println("Try asking: 'Write a haiku about programming'")

	fmt.Print("\nEnter your task: ")

	var input string
	reader := os.Stdin
	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)
	input = string(buf[:n])
	input = input[:len(input)-1]

	response, err := ag.Run(context.Background(), input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nFinal Result: %s\n", response)

	fmt.Println("\n--- Reflection Trajectory ---")
	for i, step := range ag.GetTrajectory() {
		fmt.Printf("%d. %s\n", i+1, step.Type)
		if step.Content != "" {
			fmt.Printf("   Content: %s\n", truncate(step.Content, 100))
		}
		if step.Feedback != "" {
			fmt.Printf("   Feedback: %s\n", truncate(step.Feedback, 100))
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
