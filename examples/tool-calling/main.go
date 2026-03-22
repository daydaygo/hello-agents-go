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

	customTool := &GreetingTool{}
	registry.Register(customTool)

	ag := agent.NewSimpleAgent("tool-agent", client,
		agent.WithToolRegistry(registry),
		agent.WithSystemPrompt("You are a helpful assistant with access to tools."),
	)

	fmt.Println("Tool Calling Example")
	fmt.Println("====================")
	fmt.Println("Available tools:")
	for _, name := range registry.ListTools() {
		t, _ := registry.Get(name)
		fmt.Printf("- %s: %s\n", name, t.Description())
	}

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

	fmt.Printf("\nResponse: %s\n", response)
}

type GreetingTool struct{}

func (t *GreetingTool) Name() string {
	return "greeting"
}

func (t *GreetingTool) Description() string {
	return "Generate a greeting message. Input should be a name."
}

func (t *GreetingTool) Execute(input string) (string, error) {
	return fmt.Sprintf("Hello, %s! Nice to meet you.", input), nil
}
