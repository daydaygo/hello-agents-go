package tool

import (
	"testing"
)

func TestCalculatorTool(t *testing.T) {
	tool := &CalculatorTool{}

	if tool.Name() != "calculator" {
		t.Errorf("expected name calculator, got %s", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestCalculatorExecute(t *testing.T) {
	tool := &CalculatorTool{}

	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"2 + 2", "4", false},
		{"10 - 3", "7", false},
		{"4 * 5", "20", false},
		{"15 / 3", "5", false},
		{"2 + 2 + 1", "5", false},
		{"invalid", "", true},
		{"2 / 0", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := tool.Execute(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	tools := registry.ListTools()
	if len(tools) == 0 {
		t.Error("expected at least one tool (calculator)")
	}

	tool, ok := registry.Get("calculator")
	if !ok {
		t.Fatal("expected calculator tool")
	}
	if tool.Name() != "calculator" {
		t.Errorf("expected name calculator, got %s", tool.Name())
	}
}

func TestRegistryExecute(t *testing.T) {
	registry := NewRegistry()

	result, err := registry.Execute("calculator", "2 + 2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "4" {
		t.Errorf("got %s, want 4", result)
	}
}

func TestRegistryExecuteNotFound(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Execute("nonexistent", "test")
	if err == nil {
		t.Error("expected error for nonexistent tool")
	}
}

func TestRegistryRegisterUnregister(t *testing.T) {
	registry := NewRegistry()

	custom := &testTool{name: "custom", desc: "custom tool"}
	registry.Register(custom)

	tool, ok := registry.Get("custom")
	if !ok {
		t.Fatal("expected custom tool")
	}
	if tool.Name() != "custom" {
		t.Errorf("expected name custom, got %s", tool.Name())
	}

	registry.Unregister("custom")

	_, ok = registry.Get("custom")
	if ok {
		t.Error("expected custom tool to be unregistered")
	}
}

func TestRegistryGetToolsDescription(t *testing.T) {
	registry := NewRegistry()

	desc := registry.GetToolsDescription()
	if desc == "" {
		t.Error("expected non-empty description")
	}
}

func TestRegistryGetOpenAITools(t *testing.T) {
	registry := NewRegistry()

	tools := registry.GetOpenAITools()
	if len(tools) == 0 {
		t.Error("expected at least one OpenAI tool")
	}
}

type testTool struct {
	name string
	desc string
}

func (t *testTool) Name() string        { return t.name }
func (t *testTool) Description() string { return t.desc }
func (t *testTool) Execute(input string) (string, error) {
	return "test result", nil
}
