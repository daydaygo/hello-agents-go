package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
)

// GreetingInput defines the input for the greeting tool.
type GreetingInput struct {
	Name string `json:"name" jsonschema:"required" jsonschema_description:"The name of the person to greet"`
}

func main() {
	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY environment variable")
		return
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	modelName := os.Getenv("OPENAI_MODEL")
	if modelName == "" {
		modelName = "gpt-4o"
	}

	model, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
	})
	if err != nil {
		fmt.Printf("Error creating model: %v\n", err)
		return
	}

	greetingTool, err := utils.InferTool(
		"greeting",
		"Generate a greeting message for a person by name.",
		func(ctx context.Context, input *GreetingInput) (string, error) {
			return fmt.Sprintf("Hello, %s! Nice to meet you.", input.Name), nil
		},
	)
	if err != nil {
		fmt.Printf("Error creating tool: %v\n", err)
		return
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "assistant",
		Description: "A helpful assistant with greeting capabilities",
		Instruction: "You are a helpful assistant. You can greet people using the greeting tool. Be concise and friendly.",
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{greetingTool},
			},
		},
	})
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:          agent,
		EnableStreaming: true,
	})

	fmt.Println("ChatModelAgent Example (Eino)")
	fmt.Println("==============================")
	fmt.Println("Try asking: 'Say hello to Alice' or any question")

	fmt.Print("\nEnter your question: ")

	buf := make([]byte, 1024)
	n, _ := os.Stdin.Read(buf)
	input := strings.TrimSpace(string(buf[:n]))

	fmt.Println("\n=== ChatModelAgent Response ===")

	iter := runner.Query(ctx, input)

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			fmt.Printf("Error: %v\n", event.Err)
			continue
		}

		if event.Action != nil && event.Action.TransferToAgent != nil {
			fmt.Printf("[Transfer] -> %s\n", event.Action.TransferToAgent.DestAgentName)
		}

		if msg, _, err := adk.GetMessage(event); err == nil {
			fmt.Print(msg.Content)
		}
	}

	fmt.Println("\n=== End ===")
}
