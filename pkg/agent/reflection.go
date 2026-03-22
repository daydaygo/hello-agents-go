package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/daydaygo/hello-agents-go/pkg/llm"
	"github.com/daydaygo/hello-agents-go/pkg/prompt"
)

type ReflectionAgent struct {
	name           string
	client         *llm.Client
	maxReflections int
	actPrompt      string
	reflectPrompt  string
	improvePrompt  string
	trajectory     []ReflectionStep
}

type ReflectionStep struct {
	Type     string
	Content  string
	Feedback string
}

type ReflectionAgentOption func(*ReflectionAgent)

func WithMaxReflections(max int) ReflectionAgentOption {
	return func(a *ReflectionAgent) { a.maxReflections = max }
}

func WithActPrompt(p string) ReflectionAgentOption {
	return func(a *ReflectionAgent) { a.actPrompt = p }
}

func WithReflectPrompt(p string) ReflectionAgentOption {
	return func(a *ReflectionAgent) { a.reflectPrompt = p }
}

func WithImprovePrompt(p string) ReflectionAgentOption {
	return func(a *ReflectionAgent) { a.improvePrompt = p }
}

func NewReflectionAgent(name string, client *llm.Client, opts ...ReflectionAgentOption) *ReflectionAgent {
	a := &ReflectionAgent{
		name:           name,
		client:         client,
		maxReflections: 3,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *ReflectionAgent) Run(ctx context.Context, input string) (string, error) {
	a.trajectory = nil

	result, err := a.actPhase(ctx, input)
	if err != nil {
		return "", fmt.Errorf("act phase: %w", err)
	}

	for i := 0; i < a.maxReflections; i++ {
		feedback, shouldStop, err := a.reflectPhase(ctx, input, result)
		if err != nil {
			return "", fmt.Errorf("reflect phase: %w", err)
		}

		if shouldStop {
			return result, nil
		}

		improved, err := a.improvePhase(ctx, input, result, feedback)
		if err != nil {
			return "", fmt.Errorf("improve phase: %w", err)
		}

		result = improved
	}

	return result, nil
}

func (a *ReflectionAgent) actPhase(ctx context.Context, task string) (string, error) {
	actTemplate := a.actPrompt
	if actTemplate == "" {
		tmpl, _ := prompt.Get("reflection-act")
		actTemplate = tmpl.Template
	}

	vars := map[string]string{
		"task": task,
	}

	tmpl := prompt.NewTemplate("act", actTemplate)
	systemPrompt := tmpl.Render(vars)

	result, err := a.client.Invoke(ctx, []llm.Message{
		llm.SystemMessage(systemPrompt),
		llm.UserMessage(task),
	})
	if err != nil {
		return "", fmt.Errorf("invoke act: %w", err)
	}

	a.trajectory = append(a.trajectory, ReflectionStep{
		Type:    "act",
		Content: result,
	})

	return result, nil
}

func (a *ReflectionAgent) reflectPhase(ctx context.Context, task, solution string) (string, bool, error) {
	reflectTemplate := a.reflectPrompt
	if reflectTemplate == "" {
		tmpl, _ := prompt.Get("reflection-reflect")
		reflectTemplate = tmpl.Template
	}

	vars := map[string]string{
		"task": task,
		"code": solution,
	}

	tmpl := prompt.NewTemplate("reflect", reflectTemplate)
	systemPrompt := tmpl.Render(vars)

	feedback, err := a.client.Invoke(ctx, []llm.Message{
		llm.SystemMessage(systemPrompt),
		llm.UserMessage("Analyze the solution above."),
	})
	if err != nil {
		return "", false, fmt.Errorf("invoke reflect: %w", err)
	}

	a.trajectory = append(a.trajectory, ReflectionStep{
		Type:     "reflect",
		Feedback: feedback,
	})

	shouldStop := strings.Contains(strings.ToUpper(feedback), "NO_IMPROVEMENT_NEEDED")

	return feedback, shouldStop, nil
}

func (a *ReflectionAgent) improvePhase(ctx context.Context, task, lastSolution, feedback string) (string, error) {
	improveTemplate := a.improvePrompt
	if improveTemplate == "" {
		tmpl, _ := prompt.Get("reflection-improve")
		improveTemplate = tmpl.Template
	}

	vars := map[string]string{
		"task":              task,
		"last_code_attempt": lastSolution,
		"feedback":          feedback,
	}

	tmpl := prompt.NewTemplate("improve", improveTemplate)
	systemPrompt := tmpl.Render(vars)

	result, err := a.client.Invoke(ctx, []llm.Message{
		llm.SystemMessage(systemPrompt),
		llm.UserMessage("Improve the solution based on the feedback."),
	})
	if err != nil {
		return "", fmt.Errorf("invoke improve: %w", err)
	}

	a.trajectory = append(a.trajectory, ReflectionStep{
		Type:    "improve",
		Content: result,
	})

	return result, nil
}

func (a *ReflectionAgent) GetTrajectory() []ReflectionStep {
	result := make([]ReflectionStep, len(a.trajectory))
	copy(result, a.trajectory)
	return result
}

func (a *ReflectionAgent) Name() string {
	return a.name
}
