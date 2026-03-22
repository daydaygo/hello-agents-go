package agent

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/daydaygo/hello-agents-go/pkg/llm"
	"github.com/daydaygo/hello-agents-go/pkg/prompt"
	"github.com/daydaygo/hello-agents-go/pkg/tool"
)

type ReActAgent struct {
	name         string
	client       *llm.Client
	toolRegistry *tool.Registry
	maxSteps     int
	customPrompt string
	history      []Step
}

type Step struct {
	Thought     string
	Action      string
	ActionInput string
	Observation string
}

type ReActAgentOption func(*ReActAgent)

func WithReActMaxSteps(max int) ReActAgentOption {
	return func(a *ReActAgent) { a.maxSteps = max }
}

func WithReActCustomPrompt(p string) ReActAgentOption {
	return func(a *ReActAgent) { a.customPrompt = p }
}

func NewReActAgent(name string, client *llm.Client, toolRegistry *tool.Registry, opts ...ReActAgentOption) *ReActAgent {
	a := &ReActAgent{
		name:         name,
		client:       client,
		toolRegistry: toolRegistry,
		maxSteps:     10,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *ReActAgent) Run(ctx context.Context, input string) (string, error) {
	a.history = nil

	toolsDesc := a.toolRegistry.GetToolsDescription()
	toolNames := strings.Join(a.toolRegistry.ListTools(), ", ")

	template := a.customPrompt
	if template == "" {
		tmpl, _ := prompt.Get("react")
		template = tmpl.Template
	}

	agentScratchpad := ""

	for step := 0; step < a.maxSteps; step++ {
		vars := map[string]string{
			"tools":            toolsDesc,
			"tool_names":       toolNames,
			"question":         input,
			"history":          "",
			"agent_scratchpad": agentScratchpad,
		}

		tmpl := prompt.NewTemplate("react", template)
		systemPrompt := tmpl.Render(vars)

		response, err := a.client.Invoke(ctx, []llm.Message{
			llm.SystemMessage(systemPrompt),
			llm.UserMessage(input),
		})
		if err != nil {
			return "", fmt.Errorf("react agent invoke: %w", err)
		}

		finished, answer := a.parseFinish(response)
		if finished {
			return answer, nil
		}

		thought, action, actionInput := a.parseThoughtAction(response)
		if action == "" {
			agentScratchpad += "\nThought: " + thought + "\nPlease provide a valid action."
			continue
		}

		observation, err := a.toolRegistry.Execute(action, actionInput)
		if err != nil {
			observation = fmt.Sprintf("Error: %v", err)
		}

		a.history = append(a.history, Step{
			Thought:     thought,
			Action:      action,
			ActionInput: actionInput,
			Observation: observation,
		})

		agentScratchpad += fmt.Sprintf("\nThought: %s\nAction: %s\nAction Input: %s\nObservation: %s",
			thought, action, actionInput, observation)
	}

	return "", fmt.Errorf("max steps reached without finding final answer")
}

func (a *ReActAgent) parseFinish(response string) (bool, string) {
	re := regexp.MustCompile(`(?i)Final Answer:\s*(.+)`)
	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		return true, strings.TrimSpace(matches[1])
	}
	return false, ""
}

func (a *ReActAgent) parseThoughtAction(response string) (thought, action, actionInput string) {
	thoughtRe := regexp.MustCompile(`(?i)Thought:\s*(.+?)(?:\n|$)`)
	if matches := thoughtRe.FindStringSubmatch(response); len(matches) > 1 {
		thought = strings.TrimSpace(matches[1])
	}

	actionRe := regexp.MustCompile(`(?i)Action:\s*(\w+)`)
	if matches := actionRe.FindStringSubmatch(response); len(matches) > 1 {
		action = strings.TrimSpace(matches[1])
	}

	actionInputRe := regexp.MustCompile(`(?i)Action Input:\s*(.+?)(?:\n|$)`)
	if matches := actionInputRe.FindStringSubmatch(response); len(matches) > 1 {
		actionInput = strings.TrimSpace(matches[1])
	}

	return
}

func (a *ReActAgent) GetHistory() []Step {
	result := make([]Step, len(a.history))
	copy(result, a.history)
	return result
}

func (a *ReActAgent) Name() string {
	return a.name
}
