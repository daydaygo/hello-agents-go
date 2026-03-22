package llm

import (
	"os"
	"testing"
)

func TestMessage(t *testing.T) {
	tests := []struct {
		name     string
		msg      Message
		wantRole string
		wantCnt  string
	}{
		{
			name:     "system message",
			msg:      SystemMessage("test"),
			wantRole: "system",
			wantCnt:  "test",
		},
		{
			name:     "user message",
			msg:      UserMessage("hello"),
			wantRole: "user",
			wantCnt:  "hello",
		},
		{
			name:     "assistant message",
			msg:      AssistantMessage("response"),
			wantRole: "assistant",
			wantCnt:  "response",
		},
		{
			name:     "tool result message",
			msg:      ToolResultMessage("call_123", "result"),
			wantRole: "tool",
			wantCnt:  "result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.msg.Role != tt.wantRole {
				t.Errorf("got role %s, want %s", tt.msg.Role, tt.wantRole)
			}
			if tt.msg.Content != tt.wantCnt {
				t.Errorf("got content %s, want %s", tt.msg.Content, tt.wantCnt)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	os.Unsetenv("OPENAI_MODEL")
	client := NewClient("test-key")
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.model != "gpt-4o" {
		t.Errorf("expected default model gpt-4o, got %s", client.model)
	}
}

func TestNewClientWithEnvModel(t *testing.T) {
	os.Setenv("OPENAI_MODEL", "gpt-5")
	defer os.Unsetenv("OPENAI_MODEL")

	client := NewClient("test-key")
	if client.model != "gpt-5" {
		t.Errorf("expected model gpt-5 from env, got %s", client.model)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	client := NewClient("test-key",
		WithModel("gpt-5"),
		WithTimeout(60),
		WithMaxTokens(8192),
	)

	if client.model != "gpt-5" {
		t.Errorf("expected model gpt-5, got %s", client.model)
	}
	if client.timeout != 60 {
		t.Errorf("expected timeout 60, got %v", client.timeout)
	}
	if client.maxTokens != 8192 {
		t.Errorf("expected maxTokens 8192, got %d", client.maxTokens)
	}
}

func TestConvertMessages(t *testing.T) {
	client := NewClient("test-key")

	messages := []Message{
		SystemMessage("system prompt"),
		UserMessage("user input"),
		AssistantMessage("assistant response"),
	}

	converted := client.convertMessages(messages)

	if len(converted) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(converted))
	}
}
