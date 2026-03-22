package agent

import (
	"testing"

	"github.com/daydaygo/hello-agents-go/pkg/llm"
	"github.com/daydaygo/hello-agents-go/pkg/tool"
)

func TestNewSimpleAgent(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewSimpleAgent("test", client)

	if agent.Name() != "test" {
		t.Errorf("expected name test, got %s", agent.Name())
	}
}

func TestNewSimpleAgentWithOptions(t *testing.T) {
	client := llm.NewClient("test-key")
	registry := tool.NewRegistry()

	agent := NewSimpleAgent("test", client,
		WithSystemPrompt("be helpful"),
		WithToolRegistry(registry),
		WithMaxToolIterations(10),
	)

	if agent.systemPrompt != "be helpful" {
		t.Errorf("expected system prompt, got %s", agent.systemPrompt)
	}
	if agent.toolRegistry == nil {
		t.Error("expected tool registry")
	}
	if agent.maxToolIterations != 10 {
		t.Errorf("expected max iterations 10, got %d", agent.maxToolIterations)
	}
}

func TestSimpleAgentHistory(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewSimpleAgent("test", client)

	agent.AddMessage("user", "hello")
	agent.AddMessage("assistant", "hi there")

	history := agent.GetHistory()
	if len(history) != 2 {
		t.Errorf("expected 2 messages, got %d", len(history))
	}

	agent.ClearHistory()
	history = agent.GetHistory()
	if len(history) != 0 {
		t.Errorf("expected 0 messages after clear, got %d", len(history))
	}
}

func TestSimpleAgentBuildMessages(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewSimpleAgent("test", client,
		WithSystemPrompt("system prompt"),
	)

	agent.AddMessage("user", "hello")

	messages := agent.buildMessages()

	if len(messages) != 2 {
		t.Errorf("expected 2 messages (system + user), got %d", len(messages))
	}

	if messages[0].Role != "system" {
		t.Errorf("expected first message to be system, got %s", messages[0].Role)
	}
}

func TestSimpleAgentName(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewSimpleAgent("assistant", client)

	if agent.Name() != "assistant" {
		t.Errorf("expected name assistant, got %s", agent.Name())
	}
}
