package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/daydaygo/hello-agents-go/pkg/llm"
	"github.com/daydaygo/hello-agents-go/pkg/tool"
)

type SimpleAgent struct {
	name              string
	client            *llm.Client
	systemPrompt      string
	toolRegistry      *tool.Registry
	maxToolIterations int
	history           []llm.Message
	mu                sync.RWMutex
}

type SimpleAgentOption func(*SimpleAgent)

func WithSystemPrompt(prompt string) SimpleAgentOption {
	return func(a *SimpleAgent) { a.systemPrompt = prompt }
}

func WithToolRegistry(registry *tool.Registry) SimpleAgentOption {
	return func(a *SimpleAgent) { a.toolRegistry = registry }
}

func WithMaxToolIterations(max int) SimpleAgentOption {
	return func(a *SimpleAgent) { a.maxToolIterations = max }
}

func NewSimpleAgent(name string, client *llm.Client, opts ...SimpleAgentOption) *SimpleAgent {
	a := &SimpleAgent{
		name:              name,
		client:            client,
		maxToolIterations: 5,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *SimpleAgent) Run(ctx context.Context, input string) (string, error) {
	a.mu.Lock()
	a.history = append(a.history, llm.UserMessage(input))
	a.mu.Unlock()

	messages := a.buildMessages()

	if a.toolRegistry != nil {
		return a.runWithTools(ctx, messages)
	}

	response, err := a.client.Invoke(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("simple agent run: %w", err)
	}

	a.mu.Lock()
	a.history = append(a.history, llm.AssistantMessage(response))
	a.mu.Unlock()

	return response, nil
}

func (a *SimpleAgent) runWithTools(ctx context.Context, messages []llm.Message) (string, error) {
	tools := a.toolRegistry.GetOpenAITools()

	for i := 0; i < a.maxToolIterations; i++ {
		content, toolCalls, err := a.client.InvokeWithTools(ctx, messages, tools)
		if err != nil {
			return "", fmt.Errorf("invoke with tools: %w", err)
		}

		if len(toolCalls) == 0 {
			a.mu.Lock()
			a.history = append(a.history, llm.AssistantMessage(content))
			a.mu.Unlock()
			return content, nil
		}

		messages = append(messages, llm.AssistantMessage(content))

		for _, tc := range toolCalls {
			result, err := a.toolRegistry.Execute(tc.Name, tc.Arguments)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}
			messages = append(messages, llm.ToolResultMessage(tc.ID, result))
		}
	}

	return "", fmt.Errorf("max tool iterations reached")
}

func (a *SimpleAgent) buildMessages() []llm.Message {
	var messages []llm.Message

	if a.systemPrompt != "" {
		messages = append(messages, llm.SystemMessage(a.systemPrompt))
	}

	a.mu.RLock()
	history := make([]llm.Message, len(a.history))
	copy(history, a.history)
	a.mu.RUnlock()

	messages = append(messages, history...)
	return messages
}

func (a *SimpleAgent) AddMessage(role, content string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.history = append(a.history, llm.Message{Role: role, Content: content})
}

func (a *SimpleAgent) GetHistory() []llm.Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]llm.Message, len(a.history))
	copy(result, a.history)
	return result
}

func (a *SimpleAgent) ClearHistory() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.history = nil
}

func (a *SimpleAgent) Name() string {
	return a.name
}
