package agent

import (
	"testing"

	"github.com/daydaygo/hello-agents-go/pkg/llm"
)

func TestNewPlanAndSolveAgent(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewPlanAndSolveAgent("planner", client)

	if agent.Name() != "planner" {
		t.Errorf("expected name planner, got %s", agent.Name())
	}
}

func TestNewPlanAndSolveAgentWithOptions(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewPlanAndSolveAgent("planner", client,
		WithPlannerPrompt("custom planner"),
		WithExecutorPrompt("custom executor"),
	)

	if agent.plannerPrompt != "custom planner" {
		t.Errorf("expected planner prompt, got %s", agent.plannerPrompt)
	}
	if agent.executorPrompt != "custom executor" {
		t.Errorf("expected executor prompt, got %s", agent.executorPrompt)
	}
}

func TestPlanAndSolveAgentParsePlan(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewPlanAndSolveAgent("planner", client)

	tests := []struct {
		name      string
		response  string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "json array",
			response:  `["Step 1", "Step 2", "Step 3"]`,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "numbered list",
			response:  "1. First step\n2. Second step\n3. Third step",
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "parenthesis list",
			response:  "1) First step\n2) Second step",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:     "empty response",
			response: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			steps, err := agent.parsePlan(tt.response)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(steps) != tt.wantCount {
				t.Errorf("expected %d steps, got %d", tt.wantCount, len(steps))
			}
		})
	}
}

func TestPlanAndSolveAgentGetTasks(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewPlanAndSolveAgent("planner", client)

	agent.tasks = []Task{
		{Description: "Task 1", Status: "completed", Result: "Result 1"},
		{Description: "Task 2", Status: "pending", Result: ""},
	}

	tasks := agent.GetTasks()
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Status != "completed" {
		t.Errorf("expected status completed, got %s", tasks[0].Status)
	}
}
