# hello-agents-go

A minimal Go implementation of AI Agent framework, migrated from [hello-agents](https://github.com/datawhalechina/hello-agents).

## Features

- **LLM Client**: OpenAI-compatible API client with streaming support
- **SimpleAgent**: Basic conversational agent with message history and tool calling
- **ReActAgent**: Reasoning-Action agent with Thought-Action-Observation loop
- **PlanAndSolveAgent**: Planning and execution agent for complex tasks
- **ReflectionAgent**: Self-reflection and improvement agent
- **Tool Registry**: Extensible tool system with built-in calculator
- **Prompt Templates**: Flexible prompt template management

## Requirements

- Go >= 1.25
- OpenAI API key (or compatible API)

## Installation

```bash
go get github.com/daydaygo/hello-agents-go
```

## Quick Start

### Simple Agent

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/daydaygo/hello-agents-go/pkg/agent"
    "github.com/daydaygo/hello-agents-go/pkg/llm"
)

func main() {
    client := llm.NewClient(os.Getenv("OPENAI_API_KEY"))
    ag := agent.NewSimpleAgent("assistant", client,
        agent.WithSystemPrompt("You are a helpful assistant."),
    )
    
    response, err := ag.Run(context.Background(), "Hello!")
    if err != nil {
        panic(err)
    }
    fmt.Println(response)
}
```

### ReAct Agent

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/daydaygo/hello-agents-go/pkg/agent"
    "github.com/daydaygo/hello-agents-go/pkg/llm"
    "github.com/daydaygo/hello-agents-go/pkg/tool"
)

func main() {
    client := llm.NewClient(os.Getenv("OPENAI_API_KEY"))
    registry := tool.NewRegistry()
    
    ag := agent.NewReActAgent("react", client, registry)
    
    response, err := ag.Run(context.Background(), "What is 5 + 3?")
    if err != nil {
        panic(err)
    }
    fmt.Println(response)
}
```

### Custom Tool

```go
package main

import (
    "github.com/daydaygo/hello-agents-go/pkg/tool"
)

type WeatherTool struct{}

func (t *WeatherTool) Name() string {
    return "weather"
}

func (t *WeatherTool) Description() string {
    return "Get current weather for a city"
}

func (t *WeatherTool) Execute(input string) (string, error) {
    return "Sunny, 25°C", nil
}

func main() {
    registry := tool.NewRegistry()
    registry.Register(&WeatherTool{})
}
```

## Project Structure

```
hello-agents-go/
├── pkg/
│   ├── llm/          # LLM client (OpenAI SDK wrapper)
│   ├── agent/        # Agent implementations
│   │   ├── simple.go       # SimpleAgent
│   │   ├── react.go        # ReActAgent
│   │   ├── plan_solve.go   # PlanAndSolveAgent
│   │   └── reflection.go   # ReflectionAgent
│   ├── tool/         # Tool system
│   └── prompt/       # Prompt templates
├── examples/         # Example programs
└── go.mod
```

## Configuration

Set environment variables:

```bash
export OPENAI_API_KEY=sk-xxx
export OPENAI_BASE_URL=https://api.openai.com/v1  # optional
export OPENAI_MODEL=gpt-4o                         # optional
```

## Examples

Run the example programs:

```bash
# Simple agent
go run examples/simple-agent/main.go

# ReAct agent with tools
go run examples/react-agent/main.go

# Plan-and-Solve agent
go run examples/plan-solve-agent/main.go

# Reflection agent
go run examples/reflection-agent/main.go

# Tool calling
go run examples/tool-calling/main.go
```

## API Reference

### LLM Client

```go
client := llm.NewClient(apiKey,
    llm.WithModel("gpt-5"),
    llm.WithBaseURL("https://api.openai.com/v1"),
    llm.WithTimeout(60 * time.Second),
)

// Sync call
response, err := client.Invoke(ctx, messages)

// Stream call
ch := client.StreamInvoke(ctx, messages)
for chunk := range ch {
    fmt.Print(chunk.Content)
}
```

### Agent Interface

```go
type Agent interface {
    Run(ctx context.Context, input string) (string, error)
}
```

### Tool Interface

```go
type Tool interface {
    Name() string
    Description() string
    Execute(input string) (string, error)
}
```

## License

MIT