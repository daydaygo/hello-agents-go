package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/daydaygo/hello-agents-go/pkg/llm"
	"github.com/daydaygo/hello-agents-go/pkg/prompt"
)

type Task struct {
	Description string
	Status      string
	Result      string
}

type PlanAndSolveAgent struct {
	name           string
	client         *llm.Client
	plannerPrompt  string
	executorPrompt string
	tasks          []Task
}

type PlanAndSolveAgentOption func(*PlanAndSolveAgent)

func WithPlannerPrompt(p string) PlanAndSolveAgentOption {
	return func(a *PlanAndSolveAgent) { a.plannerPrompt = p }
}

func WithExecutorPrompt(p string) PlanAndSolveAgentOption {
	return func(a *PlanAndSolveAgent) { a.executorPrompt = p }
}

func NewPlanAndSolveAgent(name string, client *llm.Client, opts ...PlanAndSolveAgentOption) *PlanAndSolveAgent {
	a := &PlanAndSolveAgent{
		name:   name,
		client: client,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *PlanAndSolveAgent) Run(ctx context.Context, input string) (string, error) {
	a.tasks = nil

	plan, err := a.planPhase(ctx, input)
	if err != nil {
		return "", fmt.Errorf("plan phase: %w", err)
	}

	result, err := a.executePhase(ctx, input, plan)
	if err != nil {
		return "", fmt.Errorf("execute phase: %w", err)
	}

	return result, nil
}

func (a *PlanAndSolveAgent) planPhase(ctx context.Context, question string) ([]string, error) {
	plannerTemplate := a.plannerPrompt
	if plannerTemplate == "" {
		tmpl, _ := prompt.Get("planner")
		plannerTemplate = tmpl.Template
	}

	vars := map[string]string{
		"question": question,
	}

	tmpl := prompt.NewTemplate("planner", plannerTemplate)
	systemPrompt := tmpl.Render(vars)

	response, err := a.client.Invoke(ctx, []llm.Message{
		llm.SystemMessage(systemPrompt),
		llm.UserMessage(question),
	})
	if err != nil {
		return nil, fmt.Errorf("invoke planner: %w", err)
	}

	steps, err := a.parsePlan(response)
	if err != nil {
		return nil, fmt.Errorf("parse plan: %w", err)
	}

	return steps, nil
}

func (a *PlanAndSolveAgent) parsePlan(response string) ([]string, error) {
	response = strings.TrimSpace(response)

	if strings.HasPrefix(response, "[") {
		var steps []string
		if err := json.Unmarshal([]byte(response), &steps); err != nil {
			return nil, fmt.Errorf("json unmarshal: %w", err)
		}
		return steps, nil
	}

	var steps []string
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if len(line) > 2 && (line[0] >= '0' && line[0] <= '9') && (line[1] == '.' || line[1] == ')') {
			step := strings.TrimSpace(line[2:])
			if step != "" {
				steps = append(steps, step)
			}
		} else if len(line) > 3 && (line[0] >= '0' && line[0] <= '9') && (line[2] >= '0' && line[2] <= '9') && line[3] == '.' {
			step := strings.TrimSpace(line[4:])
			if step != "" {
				steps = append(steps, step)
			}
		}
	}

	if len(steps) == 0 {
		return nil, fmt.Errorf("no valid steps found in response")
	}

	return steps, nil
}

func (a *PlanAndSolveAgent) executePhase(ctx context.Context, question string, steps []string) (string, error) {
	executorTemplate := a.executorPrompt
	if executorTemplate == "" {
		tmpl, _ := prompt.Get("executor")
		executorTemplate = tmpl.Template
	}

	planStr := strings.Join(steps, "\n")
	var history strings.Builder

	for i, step := range steps {
		task := Task{
			Description: step,
			Status:      "pending",
		}

		vars := map[string]string{
			"question":     question,
			"plan":         planStr,
			"history":      history.String(),
			"current_step": step,
		}

		tmpl := prompt.NewTemplate("executor", executorTemplate)
		systemPrompt := tmpl.Render(vars)

		result, err := a.client.Invoke(ctx, []llm.Message{
			llm.SystemMessage(systemPrompt),
			llm.UserMessage(step),
		})
		if err != nil {
			return "", fmt.Errorf("execute step %d: %w", i+1, err)
		}

		task.Status = "completed"
		task.Result = result
		a.tasks = append(a.tasks, task)

		fmt.Fprintf(&history, "Step %d: %s\nResult: %s\n\n", i+1, step, result)
	}

	if len(a.tasks) > 0 {
		return a.tasks[len(a.tasks)-1].Result, nil
	}

	return "", fmt.Errorf("no tasks executed")
}

func (a *PlanAndSolveAgent) GetTasks() []Task {
	result := make([]Task, len(a.tasks))
	copy(result, a.tasks)
	return result
}

func (a *PlanAndSolveAgent) Name() string {
	return a.name
}
