package main

import (
	"context"
	"fmt"
	"os"

	"github.com/daydaygo/hello-agents-go/pkg/agent"
	"github.com/daydaygo/hello-agents-go/pkg/llm"
	"github.com/daydaygo/hello-agents-go/pkg/tool"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY environment variable")
		return
	}

	client := llm.NewClient(apiKey)
	registry := tool.NewRegistry()

	ag := agent.NewReActAgent("react-agent", client, registry)

	fmt.Println("ReAct Agent Example")
	fmt.Println("===================")
	fmt.Println("Try asking: 'What is 5 + 3?' or 'Calculate 10 * 4'")

	fmt.Print("\nEnter your question: ")

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

	fmt.Printf("\nFinal Answer: %s\n", response)

	fmt.Println("\n--- Execution History ---")
	for i, step := range ag.GetHistory() {
		fmt.Printf("Step %d:\n", i+1)
		fmt.Printf("  Thought: %s\n", step.Thought)
		fmt.Printf("  Action: %s(%s)\n", step.Action, step.ActionInput)
		fmt.Printf("  Observation: %s\n", step.Observation)
	}
}
