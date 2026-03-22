package llm

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Message struct {
	Role    string
	Content string
}

type Client struct {
	client    openai.Client
	model     string
	timeout   time.Duration
	maxTokens int
}

type Option func(*Client)

func WithModel(model string) Option {
	return func(c *Client) { c.model = model }
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.client = openai.NewClient(option.WithBaseURL(url))
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) { c.timeout = timeout }
}

func WithMaxTokens(tokens int) Option {
	return func(c *Client) { c.maxTokens = tokens }
}

func NewClient(apiKey string, opts ...Option) *Client {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o"
	}

	c := &Client{
		client:    openai.NewClient(option.WithAPIKey(apiKey)),
		model:     model,
		timeout:   30 * time.Second,
		maxTokens: 4096,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Invoke(ctx context.Context, messages []Message) (string, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	params := openai.ChatCompletionNewParams{
		Model:    c.model,
		Messages: c.convertMessages(messages),
	}
	if c.maxTokens > 0 {
		params.MaxCompletionTokens = openai.Int(int64(c.maxTokens))
	}

	completion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("invoke: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return completion.Choices[0].Message.Content, nil
}

type StreamChunk struct {
	Content string
	Error   error
}

func (c *Client) StreamInvoke(ctx context.Context, messages []Message) <-chan StreamChunk {
	ch := make(chan StreamChunk, 100)

	go func() {
		defer close(ch)

		if c.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, c.timeout)
			defer cancel()
		}

		params := openai.ChatCompletionNewParams{
			Model:    c.model,
			Messages: c.convertMessages(messages),
		}
		if c.maxTokens > 0 {
			params.MaxCompletionTokens = openai.Int(int64(c.maxTokens))
		}

		stream := c.client.Chat.Completions.NewStreaming(ctx, params)
		defer func() { _ = stream.Close() }()

		for {
			if !stream.Next() {
				if err := stream.Err(); err != nil && err != io.EOF {
					ch <- StreamChunk{Error: fmt.Errorf("stream: %w", err)}
				}
				return
			}
			chunk := stream.Current()
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				ch <- StreamChunk{Content: chunk.Choices[0].Delta.Content}
			}
		}
	}()

	return ch
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}

func (c *Client) InvokeWithTools(ctx context.Context, messages []Message, tools []openai.ChatCompletionToolUnionParam) (string, []ToolCall, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	params := openai.ChatCompletionNewParams{
		Model:    c.model,
		Tools:    tools,
		Messages: c.convertMessages(messages),
	}
	if c.maxTokens > 0 {
		params.MaxCompletionTokens = openai.Int(int64(c.maxTokens))
	}

	completion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", nil, fmt.Errorf("invoke with tools: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", nil, fmt.Errorf("no response choices")
	}

	choice := completion.Choices[0]
	var toolCalls []ToolCall
	for _, tc := range choice.Message.ToolCalls {
		toolCalls = append(toolCalls, ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
		})
	}
	return choice.Message.Content, toolCalls, nil
}

func (c *Client) convertMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
	result := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, msg := range messages {
		switch msg.Role {
		case "system":
			result[i] = openai.SystemMessage(msg.Content)
		case "user":
			result[i] = openai.UserMessage(msg.Content)
		case "assistant":
			result[i] = openai.AssistantMessage(msg.Content)
		default:
			result[i] = openai.UserMessage(msg.Content)
		}
	}
	return result
}

func SystemMessage(content string) Message {
	return Message{Role: "system", Content: content}
}

func UserMessage(content string) Message {
	return Message{Role: "user", Content: content}
}

func AssistantMessage(content string) Message {
	return Message{Role: "assistant", Content: content}
}

func ToolResultMessage(toolCallID, content string) Message {
	return Message{Role: "tool", Content: content}
}
