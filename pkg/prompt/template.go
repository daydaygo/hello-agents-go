package prompt

import (
	"fmt"
	"strings"
	"sync"
)

type Template struct {
	Name     string
	Template string
}

func NewTemplate(name, template string) *Template {
	return &Template{
		Name:     name,
		Template: template,
	}
}

func (t *Template) Render(vars map[string]string) string {
	result := t.Template
	for key, value := range vars {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}

type Registry struct {
	mu        sync.RWMutex
	templates map[string]*Template
}

func NewRegistry() *Registry {
	r := &Registry{
		templates: make(map[string]*Template),
	}
	r.registerDefaults()
	return r
}

func (r *Registry) Register(t *Template) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[t.Name] = t
}

func (r *Registry) Get(name string) (*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[name]
	if !ok {
		return nil, fmt.Errorf("template %q not found", name)
	}
	return t, nil
}

func (r *Registry) Render(name string, vars map[string]string) (string, error) {
	t, err := r.Get(name)
	if err != nil {
		return "", err
	}
	return t.Render(vars), nil
}

func (r *Registry) registerDefaults() {
	r.Register(NewTemplate("simple-agent", `You are a helpful assistant.

{{if .SystemPrompt}}System Instructions:
{{.SystemPrompt}}

{{end}}User: {{.Input}}
Assistant:`))

	r.Register(NewTemplate("react", `You are a helpful assistant that can use tools to solve problems.

Available tools:
{tools}

Use the following format:

Thought: your reasoning about what to do
Action: the action to take, should be one of [{tool_names}]
Action Input: the input to the action
Observation: the result of the action
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: I now know the final answer
Final Answer: the final answer to the original input question

Begin!

Question: {question}
{history}
Thought: {agent_scratchpad}`))

	r.Register(NewTemplate("planner", `You are a planning assistant. Given a complex question, break it down into a list of executable steps.

Question: {question}

Provide a numbered list of steps to solve this question. Each step should be a clear, actionable task.
Format your response as a JSON array of strings, e.g.: ["Step 1", "Step 2", "Step 3"]`))

	r.Register(NewTemplate("executor", `You are an execution assistant. Execute the given step and provide the result.

Original Question: {question}
Overall Plan: {plan}
Previous Steps and Results: {history}

Current Step: {current_step}

Execute this step and provide the result.`))

	r.Register(NewTemplate("reflection-act", `You are a helpful assistant. Complete the given task.

Task: {task}

Provide your solution.`))

	r.Register(NewTemplate("reflection-reflect", `You are a critical reviewer. Analyze the following solution and identify any issues or areas for improvement.

Task: {task}
Solution: {code}

If the solution is perfect and needs no improvement, respond with "NO_IMPROVEMENT_NEEDED".
Otherwise, list specific issues and suggestions for improvement.`))

	r.Register(NewTemplate("reflection-improve", `You are a helpful assistant. Improve the solution based on the feedback.

Task: {task}
Previous Solution: {last_code_attempt}
Feedback: {feedback}

Provide an improved solution.`))
}

var defaultRegistry = NewRegistry()

func Register(t *Template) {
	defaultRegistry.Register(t)
}

func Get(name string) (*Template, error) {
	return defaultRegistry.Get(name)
}

func Render(name string, vars map[string]string) (string, error) {
	return defaultRegistry.Render(name, vars)
}
