package agent

import (
	"testing"

	"github.com/daydaygo/hello-agents-go/pkg/llm"
)

func TestNewReflectionAgent(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewReflectionAgent("reflector", client)

	if agent.Name() != "reflector" {
		t.Errorf("expected name reflector, got %s", agent.Name())
	}
	if agent.maxReflections != 3 {
		t.Errorf("expected default max reflections 3, got %d", agent.maxReflections)
	}
}

func TestNewReflectionAgentWithOptions(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewReflectionAgent("reflector", client,
		WithMaxReflections(5),
		WithActPrompt("custom act"),
		WithReflectPrompt("custom reflect"),
		WithImprovePrompt("custom improve"),
	)

	if agent.maxReflections != 5 {
		t.Errorf("expected max reflections 5, got %d", agent.maxReflections)
	}
	if agent.actPrompt != "custom act" {
		t.Errorf("expected act prompt, got %s", agent.actPrompt)
	}
	if agent.reflectPrompt != "custom reflect" {
		t.Errorf("expected reflect prompt, got %s", agent.reflectPrompt)
	}
	if agent.improvePrompt != "custom improve" {
		t.Errorf("expected improve prompt, got %s", agent.improvePrompt)
	}
}

func TestReflectionAgentGetTrajectory(t *testing.T) {
	client := llm.NewClient("test-key")
	agent := NewReflectionAgent("reflector", client)

	agent.trajectory = []ReflectionStep{
		{Type: "act", Content: "initial solution"},
		{Type: "reflect", Feedback: "needs improvement"},
		{Type: "improve", Content: "improved solution"},
	}

	trajectory := agent.GetTrajectory()
	if len(trajectory) != 3 {
		t.Errorf("expected 3 steps, got %d", len(trajectory))
	}
	if trajectory[0].Type != "act" {
		t.Errorf("expected type act, got %s", trajectory[0].Type)
	}
	if trajectory[1].Type != "reflect" {
		t.Errorf("expected type reflect, got %s", trajectory[1].Type)
	}
	if trajectory[2].Type != "improve" {
		t.Errorf("expected type improve, got %s", trajectory[2].Type)
	}
}

func TestReflectionStep(t *testing.T) {
	step := ReflectionStep{
		Type:     "act",
		Content:  "test content",
		Feedback: "test feedback",
	}

	if step.Type != "act" {
		t.Errorf("expected type act, got %s", step.Type)
	}
	if step.Content != "test content" {
		t.Errorf("expected content, got %s", step.Content)
	}
	if step.Feedback != "test feedback" {
		t.Errorf("expected feedback, got %s", step.Feedback)
	}
}
