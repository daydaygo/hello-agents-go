package tool

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *Registry {
	r := &Registry{
		tools: make(map[string]Tool),
	}
	r.Register(&CalculatorTool{})
	return r
}

func (r *Registry) Register(t Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[t.Name()] = t
}

func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tools, name)
}

func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

func (r *Registry) Execute(name, input string) (string, error) {
	t, ok := r.Get(name)
	if !ok {
		return "", fmt.Errorf("tool %q not found", name)
	}
	return t.Execute(input)
}

func (r *Registry) ListTools() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

func (r *Registry) GetToolsDescription() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.tools) == 0 {
		return "No tools available."
	}
	var sb strings.Builder
	sb.WriteString("Available tools:\n")
	for name, t := range r.tools {
		fmt.Fprintf(&sb, "- %s: %s\n", name, t.Description())
	}
	return sb.String()
}

func (r *Registry) GetOpenAITools() []openai.ChatCompletionToolUnionParam {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tools := make([]openai.ChatCompletionToolUnionParam, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
			Name:        t.Name(),
			Description: openai.String(t.Description()),
			Parameters:  shared.FunctionParameters{},
		}))
	}
	return tools
}

type CalculatorTool struct{}

func (t *CalculatorTool) Name() string {
	return "calculator"
}

func (t *CalculatorTool) Description() string {
	return "Evaluate a mathematical expression. Input should be a valid mathematical expression like '2 + 2' or '3 * 4'."
}

func (t *CalculatorTool) Execute(input string) (string, error) {
	input = strings.TrimSpace(input)
	tokens := strings.Fields(input)
	if len(tokens) < 3 {
		return "", fmt.Errorf("invalid expression: need at least 3 tokens (e.g., '2 + 2')")
	}

	var result float64
	var err error

	result, err = strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", tokens[0])
	}

	for i := 1; i < len(tokens); i += 2 {
		if i+1 >= len(tokens) {
			break
		}
		op := tokens[i]
		nextNum, err := strconv.ParseFloat(tokens[i+1], 64)
		if err != nil {
			return "", fmt.Errorf("invalid number: %s", tokens[i+1])
		}

		switch op {
		case "+":
			result += nextNum
		case "-":
			result -= nextNum
		case "*":
			result *= nextNum
		case "/":
			if nextNum == 0 {
				return "", fmt.Errorf("division by zero")
			}
			result /= nextNum
		default:
			return "", fmt.Errorf("unknown operator: %s", op)
		}
	}

	if result == float64(int(result)) {
		return fmt.Sprintf("%d", int(result)), nil
	}
	return fmt.Sprintf("%.2f", result), nil
}

var defaultRegistry = NewRegistry()

func Register(t Tool) {
	defaultRegistry.Register(t)
}

func Unregister(name string) {
	defaultRegistry.Unregister(name)
}

func Get(name string) (Tool, bool) {
	return defaultRegistry.Get(name)
}

func Execute(name, input string) (string, error) {
	return defaultRegistry.Execute(name, input)
}

func ListTools() []string {
	return defaultRegistry.ListTools()
}

func GetToolsDescription() string {
	return defaultRegistry.GetToolsDescription()
}

func GetOpenAITools() []openai.ChatCompletionToolUnionParam {
	return defaultRegistry.GetOpenAITools()
}
