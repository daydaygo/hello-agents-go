package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
)

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

	// Create calculator tools
	addTool, err := utils.InferTool(
		"add",
		"Add two numbers together.",
		func(ctx context.Context, input *struct {
			A float64 `json:"a" jsonschema:"required" jsonschema_description:"First number"`
			B float64 `json:"b" jsonschema:"required" jsonschema_description:"Second number"`
		}) (string, error) {
			result := input.A + input.B
			if result == float64(int(result)) {
				return fmt.Sprintf("%d", int(result)), nil
			}
			return fmt.Sprintf("%.2f", result), nil
		},
	)
	if err != nil {
		fmt.Printf("Error creating add tool: %v\n", err)
		return
	}

	multiplyTool, err := utils.InferTool(
		"multiply",
		"Multiply two numbers together.",
		func(ctx context.Context, input *struct {
			A float64 `json:"a" jsonschema:"required" jsonschema_description:"First number"`
			B float64 `json:"b" jsonschema:"required" jsonschema_description:"Second number"`
		}) (string, error) {
			result := input.A * input.B
			if result == float64(int(result)) {
				return fmt.Sprintf("%d", int(result)), nil
			}
			return fmt.Sprintf("%.2f", result), nil
		},
	)
	if err != nil {
		fmt.Printf("Error creating multiply tool: %v\n", err)
		return
	}

	expressionTool, err := utils.InferTool(
		"evaluate_expression",
		"Evaluate a math expression string like '2 + 3' or '10 * 4 - 1'. Supports +, -, *, /.",
		func(ctx context.Context, input *struct {
			Expression string `json:"expression" jsonschema:"required" jsonschema_description:"Math expression to evaluate"`
		}) (string, error) {
			return calculate(input.Expression)
		},
	)
	if err != nil {
		fmt.Printf("Error creating expression tool: %v\n", err)
		return
	}

	// Create sub-agents
	mathAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "math_agent",
		Description: "Perform mathematical calculations including addition, multiplication, and expression evaluation",
		Instruction: "You are a math specialist. Use the available tools to perform calculations accurately.",
		Model:       model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{addTool, multiplyTool, expressionTool},
			},
		},
	})
	if err != nil {
		fmt.Printf("Error creating math agent: %v\n", err)
		return
	}

	generalAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "general_agent",
		Description: "Handle general knowledge questions and provide information",
		Instruction: "You are a general knowledge specialist. Help users find information and answer questions.",
		Model:       model,
	})
	if err != nil {
		fmt.Printf("Error creating general agent: %v\n", err)
		return
	}

	// Create the supervisor
	supervisorAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "supervisor",
		Description: "Route tasks to the appropriate specialist agent",
		Instruction: "You are a supervisor. Route math questions to math_agent and general questions to general_agent. Always delegate to a specialist.",
		Model:       model,
	})
	if err != nil {
		fmt.Printf("Error creating supervisor: %v\n", err)
		return
	}

	// Assemble supervisor pattern
	agent, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: supervisorAgent,
		SubAgents:  []adk.Agent{mathAgent, generalAgent},
	})
	if err != nil {
		fmt.Printf("Error creating supervisor pattern: %v\n", err)
		return
	}

	// Run
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:          agent,
		EnableStreaming: true,
	})

	fmt.Println("Supervisor Multi-Agent Example (Eino)")
	fmt.Println("======================================")
	fmt.Println("Try asking: 'What is 15 + 27? And what is 8 * 9?'")

	fmt.Print("\nEnter your question: ")

	buf := make([]byte, 1024)
	n, _ := os.Stdin.Read(buf)
	input := strings.TrimSpace(string(buf[:n]))

	fmt.Println("\n=== Supervisor Response ===")

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

func calculate(expression string) (string, error) {
	expression = strings.TrimSpace(expression)
	parts := strings.Fields(expression)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty expression")
	}

	result, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", parts[0])
	}

	for i := 1; i+1 < len(parts); i += 2 {
		op := parts[i]
		num, err := strconv.ParseFloat(parts[i+1], 64)
		if err != nil {
			return "", fmt.Errorf("invalid number: %s", parts[i+1])
		}
		switch op {
		case "+":
			result += num
		case "-":
			result -= num
		case "*":
			result *= num
		case "/":
			if num == 0 {
				return "", fmt.Errorf("division by zero")
			}
			result /= num
		default:
			return "", fmt.Errorf("unsupported operator: %s", op)
		}
	}

	if result == float64(int(result)) {
		return fmt.Sprintf("%d", int(result)), nil
	}
	return fmt.Sprintf("%.2f", result), nil
}
