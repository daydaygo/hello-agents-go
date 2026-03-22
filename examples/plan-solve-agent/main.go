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
	ag := agent.NewPlanAndSolveAgent("planner", client)

	fmt.Println("Plan-and-Solve Agent Example")
	fmt.Println("============================")
	fmt.Println("Try asking a multi-step question like:")
	fmt.Println("'What are the differences between Go and Python for web development?'")

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

	fmt.Printf("\nFinal Result: %s\n", response)

	fmt.Println("\n--- Tasks ---")
	for i, task := range ag.GetTasks() {
		fmt.Printf("%d. %s\n", i+1, task.Description)
		fmt.Printf("   Status: %s\n", task.Status)
		fmt.Printf("   Result: %s\n", task.Result)
	}
}
