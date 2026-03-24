package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// CalculatorInput defines the input for the calculator tool.
type CalculatorInput struct {
	Expression string `json:"expression" jsonschema:"required" jsonschema_description:"Math expression like '2 + 3' or '10 * 4'. Supports +, -, *, /."`
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

	calculatorTool, err := utils.InferTool(
		"calculator",
		"Calculate a math expression. Supports basic arithmetic: +, -, *, /.",
		func(ctx context.Context, input *CalculatorInput) (string, error) {
			return calculate(input.Expression)
		},
	)
	if err != nil {
		fmt.Printf("Error creating tool: %v\n", err)
		return
	}

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: model,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{calculatorTool},
		},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			res := make([]*schema.Message, 0, len(input)+1)
			res = append(res, schema.SystemMessage("You are a helpful assistant with access to a calculator tool. Use it when asked math questions."))
			res = append(res, input...)
			return res
		},
		MaxStep: 20,
	})
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	fmt.Println("ReAct Agent Example (Eino)")
	fmt.Println("==========================")
	fmt.Println("Try asking: 'What is 5 + 3?' or 'Calculate 10 * 4'")

	fmt.Print("\nEnter your question: ")

	buf := make([]byte, 1024)
	n, _ := os.Stdin.Read(buf)
	input := strings.TrimSpace(string(buf[:n]))

	fmt.Println("\n=== ReAct Agent Response ===")

	sr, err := agent.Stream(ctx, []*schema.Message{
		schema.UserMessage(input),
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer sr.Close()

	for {
		msg, err := sr.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Print(msg.Content)
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
