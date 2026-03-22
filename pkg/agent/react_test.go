package agent

import (
	"testing"

	"github.com/daydaygo/hello-agents-go/pkg/llm"
	"github.com/daydaygo/hello-agents-go/pkg/tool"
)

func TestNewReActAgent(t *testing.T) {
	client := llm.NewClient("test-key")
	registry := tool.NewRegistry()

	agent := NewReActAgent("react", client, registry)

	if agent.Name() != "react" {
		t.Errorf("expected name react, got %s", agent.Name())
	}
	if agent.maxSteps != 10 {
		t.Errorf("expected default max steps 10, got %d", agent.maxSteps)
	}
}

func TestNewReActAgentWithOptions(t *testing.T) {
	client := llm.NewClient("test-key")
	registry := tool.NewRegistry()

	agent := NewReActAgent("react", client, registry,
		WithReActMaxSteps(5),
		WithReActCustomPrompt("custom prompt"),
	)

	if agent.maxSteps != 5 {
		t.Errorf("expected max steps 5, got %d", agent.maxSteps)
	}
	if agent.customPrompt != "custom prompt" {
		t.Errorf("expected custom prompt, got %s", agent.customPrompt)
	}
}

func TestReActAgentParseFinish(t *testing.T) {
	client := llm.NewClient("test-key")
	registry := tool.NewRegistry()
	agent := NewReActAgent("react", client, registry)

	tests := []struct {
		response   string
		wantFinish bool
		wantAnswer string
	}{
		{
			response:   "Thought: I know the answer\nFinal Answer: 42",
			wantFinish: true,
			wantAnswer: "42",
		},
		{
			response:   "Final Answer: The result is 100",
			wantFinish: true,
			wantAnswer: "The result is 100",
		},
		{
			response:   "Thought: Need more info\nAction: calculator\nAction Input: 1+1",
			wantFinish: false,
			wantAnswer: "",
		},
		{
			response:   "final answer: lowercase test",
			wantFinish: true,
			wantAnswer: "lowercase test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.response, func(t *testing.T) {
			finished, answer := agent.parseFinish(tt.response)
			if finished != tt.wantFinish {
				t.Errorf("expected finished %v, got %v", tt.wantFinish, finished)
			}
			if finished && answer != tt.wantAnswer {
				t.Errorf("expected answer %q, got %q", tt.wantAnswer, answer)
			}
		})
	}
}

func TestReActAgentParseThoughtAction(t *testing.T) {
	client := llm.NewClient("test-key")
	registry := tool.NewRegistry()
	agent := NewReActAgent("react", client, registry)

	response := "Thought: I need to calculate\nAction: calculator\nAction Input: 2 + 2"
	thought, action, actionInput := agent.parseThoughtAction(response)

	if thought != "I need to calculate" {
		t.Errorf("expected thought, got %s", thought)
	}
	if action != "calculator" {
		t.Errorf("expected action calculator, got %s", action)
	}
	if actionInput != "2 + 2" {
		t.Errorf("expected action input, got %s", actionInput)
	}
}

func TestReActAgentGetHistory(t *testing.T) {
	client := llm.NewClient("test-key")
	registry := tool.NewRegistry()
	agent := NewReActAgent("react", client, registry)

	agent.history = []Step{
		{Thought: "thinking", Action: "calculator", ActionInput: "1+1", Observation: "2"},
	}

	history := agent.GetHistory()
	if len(history) != 1 {
		t.Errorf("expected 1 step, got %d", len(history))
	}
	if history[0].Thought != "thinking" {
		t.Errorf("expected thought, got %s", history[0].Thought)
	}
}
