package prompt

import (
	"testing"
)

func TestTemplateRender(t *testing.T) {
	tests := []struct {
		name     string
		template string
		vars     map[string]string
		want     string
	}{
		{
			name:     "single variable",
			template: "Hello, {name}!",
			vars:     map[string]string{"name": "World"},
			want:     "Hello, World!",
		},
		{
			name:     "multiple variables",
			template: "{greeting}, {name}!",
			vars:     map[string]string{"greeting": "Hi", "name": "Alice"},
			want:     "Hi, Alice!",
		},
		{
			name:     "no variables",
			template: "Hello World!",
			vars:     map[string]string{},
			want:     "Hello World!",
		},
		{
			name:     "missing variable preserved",
			template: "Hello, {name}!",
			vars:     map[string]string{},
			want:     "Hello, {name}!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := NewTemplate("test", tt.template)
			got := tmpl.Render(tt.vars)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	tmpl, err := registry.Get("react")
	if err != nil {
		t.Fatalf("expected react template, got error: %v", err)
	}
	if tmpl.Name != "react" {
		t.Errorf("expected name react, got %s", tmpl.Name)
	}

	if _, err = registry.Get("planner"); err != nil {
		t.Fatalf("expected planner template, got error: %v", err)
	}

	if _, err = registry.Get("executor"); err != nil {
		t.Fatalf("expected executor template, got error: %v", err)
	}

	if _, err = registry.Get("reflection-act"); err != nil {
		t.Fatalf("expected reflection-act template, got error: %v", err)
	}

	if _, err = registry.Get("reflection-reflect"); err != nil {
		t.Fatalf("expected reflection-reflect template, got error: %v", err)
	}

	if _, err = registry.Get("reflection-improve"); err != nil {
		t.Fatalf("expected reflection-improve template, got error: %v", err)
	}
}

func TestRegistryNotFound(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent template")
	}
}

func TestRegistryRegister(t *testing.T) {
	registry := NewRegistry()

	custom := NewTemplate("custom", "Custom: {value}")
	registry.Register(custom)

	tmpl, err := registry.Get("custom")
	if err != nil {
		t.Fatalf("expected custom template, got error: %v", err)
	}
	if tmpl.Name != "custom" {
		t.Errorf("expected name custom, got %s", tmpl.Name)
	}
}

func TestRegistryRender(t *testing.T) {
	registry := NewRegistry()

	result, err := registry.Render("react", map[string]string{
		"tools":            "calculator",
		"tool_names":       "calculator",
		"question":         "test question",
		"history":          "",
		"agent_scratchpad": "",
	})
	if err != nil {
		t.Fatalf("expected render success, got error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestDefaultRegistry(t *testing.T) {
	tmpl, err := Get("react")
	if err != nil {
		t.Fatalf("expected react template, got error: %v", err)
	}
	if tmpl.Name != "react" {
		t.Errorf("expected name react, got %s", tmpl.Name)
	}

	result, err := Render("planner", map[string]string{
		"question": "test",
	})
	if err != nil {
		t.Fatalf("expected render success, got error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result")
	}
}
